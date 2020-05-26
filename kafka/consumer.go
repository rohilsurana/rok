package kafka

import (
	"log"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Consumer object for storing kafka producer, topic and server information
type Consumer struct {
	consumer *kafka.Consumer
	Topic    string
	Brokers  string
	Configs  map[string]string
}

// NewConsumer creates a new consumer for given topic and configs
func NewConsumer(brokers, topic string, configs map[string]string) *Consumer {
	log.Print("Creating consumer with config - ", brokers, " and ", topic)
	consumer, err := kafka.NewConsumer(getKafkaConfigMap(brokers, configs))
	if err != nil {
		log.Fatal("Failed to create consumer: ", err)
	}

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatal("Failed to create consumer: ", err)
	}

	return &Consumer{
		consumer: consumer,
		Topic:    topic,
		Configs:  configs,
		Brokers:  brokers,
	}
}

// Consume function sets the processing function for each Kafka Message
func (c *Consumer) Consume(processMessage func(msg *Message) error) error {
	kafkaMessage, err := c.consumer.ReadMessage(-1)
	if err != nil {
		return err
	}
	err = processMessage(&Message{kafkaMessage: kafkaMessage})

	if err != nil {
		return err
	}

	_, err = c.consumer.CommitMessage(kafkaMessage)
	return err
}
