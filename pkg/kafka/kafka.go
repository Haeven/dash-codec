package kafka

import (
	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	Reader *kafka.Reader
	Writer *kafka.Writer
}

func NewKafkaClient(brokers []string, groupID, inputTopic, outputTopic string) *KafkaClient {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    inputTopic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    outputTopic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaClient{
		Reader: reader,
		Writer: writer,
	}
}

func (c *KafkaClient) Close() error {
	if err := c.Reader.Close(); err != nil {
		return err
	}
	return c.Writer.Close()
}
