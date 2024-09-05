package kafka

import (
	"context"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaClient struct {
	Client      *kgo.Client
	InputTopic  string
	OutputTopic string
}

func NewKafkaClient(brokers []string, groupID, inputTopic, outputTopic string) (*KafkaClient, error) {
	log.Printf("Creating new Kafka client with brokers: %v, groupID: %s, inputTopic: %s, outputTopic: %s", brokers, groupID, inputTopic, outputTopic)
	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(inputTopic),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		log.Printf("Error creating Kafka client: %v", err)
		return nil, err
	}

	log.Println("Kafka client created successfully")
	return &KafkaClient{
		Client:      client,
		InputTopic:  inputTopic,
		OutputTopic: outputTopic,
	}, nil
}

func (c *KafkaClient) ReadMessage(ctx context.Context) (*kgo.Record, error) {
	log.Println("Attempting to read message from Kafka")
	fetches := c.Client.PollFetches(ctx)
	if errs := fetches.Errors(); len(errs) > 0 {
		log.Printf("Error polling fetches: %v", errs[0].Err)
		return nil, errs[0].Err
	}

	records := fetches.Records()
	if len(records) > 0 {
		log.Println("Message read successfully from Kafka")
		return records[0], nil
	}

	log.Println("No messages available to read from Kafka")
	return nil, nil
}

func (c *KafkaClient) WriteMessage(ctx context.Context, key, value []byte) error {
	log.Println("Attempting to write message to Kafka")
	record := &kgo.Record{
		Topic: c.OutputTopic,
		Key:   key,
		Value: value,
	}
	err := c.Client.ProduceSync(ctx, record).FirstErr()
	if err != nil {
		log.Printf("Error writing message to Kafka: %v", err)
		return err
	}
	log.Println("Message written successfully to Kafka")
	return nil
}

func (c *KafkaClient) Close() error {
	log.Println("Closing Kafka client")
	c.Client.Close()
	log.Println("Kafka client closed")
	return nil
}
