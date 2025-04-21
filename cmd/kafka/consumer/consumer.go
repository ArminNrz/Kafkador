package consumer

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"kafkador/cmd/kafka/model"
	"kafkador/cmd/service"
	"log"
	_ "time"
)

type KafkaConsumer struct {
	Reader       *kafka.Reader
	BikerService service.BikerService
}

func NewKafkaConsumer(broker string, topic string, groupID string, bikerService service.BikerService) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return &KafkaConsumer{
		Reader:       reader,
		BikerService: bikerService,
	}
}

func (consumer *KafkaConsumer) Consume(ctx context.Context) {
	for {
		m, err := consumer.Reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("error while reading message: %v", err)
			continue
		}

		var msg model.CreateBikerResponseMessage
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			continue
		}

		log.Printf("Received message: %+v", msg)
		consumer.BikerService.HandleBikerResponse(msg)
	}
}
