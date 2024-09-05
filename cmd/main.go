package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"codec/internal/server"
	"codec/pkg/kafka"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	brokers := []string{"localhost:9092"} // Replace with your Kafka brokers
	groupID := "video-processor"
	inputTopic := "video-uploads"
	outputTopic := "processed-videos"

	kafkaClient := kafka.NewKafkaClient(brokers, groupID, inputTopic, outputTopic)
	defer kafkaClient.Close()

	srv := server.NewServer(kafkaClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		log.Println("Received shutdown signal. Shutting down...")
		cancel()
	}()

	return srv.Run(ctx)
}
