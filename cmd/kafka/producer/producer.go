package producer

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"kafkador/cmd/kafka/model"
	"log"
	"time"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

func NewKafkaProducer(broker string, topic string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{
		Writer: writer,
	}
}

func (producer *KafkaProducer) Produce(ctx context.Context, message model.CreateBikerMessage) {
	value, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("json marshal error: %v", err)
	}

	err = producer.Writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(message.CorrelationID),
		Value: value,
		Time:  time.Now(),
	})
	if err != nil {
		log.Fatalf("failed to write message: %v", err)
	}
}
