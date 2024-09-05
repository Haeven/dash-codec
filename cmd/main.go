package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Haeven/codec/internal/server"
	"github.com/Haeven/codec/pkg/kafka"
)

func main() {
	log.Println("Starting application")
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Println("Application shutdown complete")
}

func run() error {
	brokers := []string{"localhost:9092"}
	groupID := "video-processor"
	inputTopic := "video-uploads"
	outputTopic := "processed-videos"

	log.Println("Initializing Kafka client")
	kafkaClient, err := kafka.NewKafkaClient(brokers, groupID, inputTopic, outputTopic)
	if err != nil {
		log.Printf("Failed to create Kafka client: %v", err)
		return err
	}
	defer func() {
		log.Println("Closing Kafka client")
		kafkaClient.Close()
	}()

	log.Println("Creating server instance")
	srv := server.NewServer(kafkaClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		sig := <-sigCh
		log.Printf("Received shutdown signal: %v. Shutting down...", sig)
		cancel()
	}()

	log.Println("Starting server")
	err = srv.Run(ctx)
	if err != nil {
		log.Printf("Server stopped with error: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
	return err
}
