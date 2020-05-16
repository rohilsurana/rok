package kafka

import (
	"log"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Producer object for storing kafka producer, topic and server information
type Producer struct {
	producer         *kafka.Producer
	Topic            string
	BootstrapServers string
}

// NewProducer creates a new producer for given bootstrap servers and topic
func NewProducer(bootstrapServers, topic string) *Producer {
	log.Print("Creating producer with config - ", bootstrapServers, " and ", topic)
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})
	if err != nil {
		log.Fatal("Failed to create producer: ", err)
	}
	return &Producer{
		producer:         p,
		Topic:            topic,
		BootstrapServers: bootstrapServers,
	}
}

// SendSync sends a message using the intialised producer
func (p *Producer) SendSync(message []byte) error {
	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.Topic,
			Partition: kafka.PartitionAny,
		},
		Value: message,
	}, nil)
	if err != nil {
		return err
	}
	deliveryEvent := (<-p.producer.Events()).(*kafka.Message)
	if deliveryEvent.TopicPartition.Error != nil {
		return deliveryEvent.TopicPartition.Error
	}
	return nil
}
