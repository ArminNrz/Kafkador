package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"kafkador/cmd/api"
	"kafkador/cmd/kafka/model"
	"kafkador/cmd/kafka/producer"
	"log"
	"sync"
	"time"
)

type BikerService interface {
	CreateBiker(ctx context.Context, request api.CreateBikerRequest) (api.CreateBikerResponse, error)
	HandleBikerResponse(msg model.CreateBikerResponseMessage)
	StartCleanupLoop(interval time.Duration, timeout time.Duration)
}

type BikerServiceImpl struct {
	Kafka         *producer.KafkaProducer
	ResponseChans sync.Map
}

type ResponseWrapper struct {
	Chan      chan model.CreateBikerResponseMessage
	CreatedAt time.Time
}

func NewBikerService(kafka *producer.KafkaProducer) BikerService {
	return &BikerServiceImpl{
		Kafka: kafka,
	}
}

func (b *BikerServiceImpl) CreateBiker(ctx context.Context, request api.CreateBikerRequest) (api.CreateBikerResponse, error) {
	correlationID := uuid.New().String()
	message := model.CreateBikerMessage{
		CorrelationID: correlationID,
		Payload:       request,
	}

	wrapper := &ResponseWrapper{
		Chan:      make(chan model.CreateBikerResponseMessage, 1),
		CreatedAt: time.Now(),
	}
	b.ResponseChans.Store(correlationID, wrapper)
	defer b.ResponseChans.Delete(correlationID)

	go b.Kafka.Produce(ctx, message)

	select {
	case response := <-wrapper.Chan:
		if response.StatusCode != 200 {
			return api.CreateBikerResponse{}, errors.New(response.Message)
		}
		return api.CreateBikerResponse{
			ID:          response.ID,
			Name:        response.Name,
			PhoneNumber: response.PhoneNumber,
		}, nil
	case <-time.After(10 * time.Second):
		return api.CreateBikerResponse{}, errors.New("timeout waiting for biker creation response")
	}
}

func (b *BikerServiceImpl) HandleBikerResponse(msg model.CreateBikerResponseMessage) {
	value, ok := b.ResponseChans.Load(msg.CorrelationID)
	if !ok {
		// No one is waiting for this correlation ID
		log.Printf("[WARN] No waiting channel for CorrelationID: %s", msg.CorrelationID)
		return
	}

	wrapper, ok := value.(*ResponseWrapper)
	if !ok {
		log.Printf("[ERROR] Invalid type assertion for CorrelationID: %s", msg.CorrelationID)
		return
	}

	select {
	case wrapper.Chan <- msg:
		// Successfully sent the message to waiting goroutine
		log.Printf("[INFO] Response sent to waiting client: %+v", msg)
	default:
		// Channel buffer is full — probably already timed out
		log.Printf("[WARN] Dropping late response for CorrelationID: %s", msg.CorrelationID)
	}
}

func (b *BikerServiceImpl) StartCleanupLoop(interval time.Duration, timeout time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now()
			b.ResponseChans.Range(func(key, value any) bool {
				wrapper, ok := value.(*ResponseWrapper)
				if !ok {
					// If not the expected type, clean it up anyway
					b.ResponseChans.Delete(key)
					return true
				}

				if now.Sub(wrapper.CreatedAt) > timeout {
					log.Printf("[CLEANUP] Removing expired correlationID: %s", key)
					close(wrapper.Chan) // Safe to close if writer is done
					b.ResponseChans.Delete(key)
				}
				return true
			})
		}
	}()
}
