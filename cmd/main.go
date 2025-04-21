package main

import (
	"context"
	"kafkador/cmd/http"
	"kafkador/cmd/kafka/consumer"
	"kafkador/cmd/kafka/producer"
	"kafkador/cmd/service"
	util "kafkador/cmd/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	producer := producer.NewKafkaProducer(config.KafkaBroker, config.ProducerTopic)

	bikerService := service.NewBikerService(producer)

	consumer := consumer.NewKafkaConsumer(
		config.KafkaBroker,
		config.ConsumerTopic,
		config.ConsumerGroupID,
		bikerService,
	)
	ctx, cancel := context.WithCancel(context.Background())
	go consumer.Consume(ctx)

	_, err = http.NewServer(config, bikerService)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	log.Println("Server is running on port:", config.HTTPServerAddress)

	// 6. Graceful shutdown
	gracefulShutdown(cancel)
}

func gracefulShutdown(cancel context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server and consumer...")
	cancel()
}
