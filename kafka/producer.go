package kafka

import (
	"errors"
	"log"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Producer object for storing kafka producer, topic and server information
type Producer struct {
	producer *kafka.Producer
	Topic    string
	Brokers  string
	Configs  map[string]string
}

// NewProducer creates a new producer for given bootstrap servers and topic
func NewProducer(brokers, topic string, configs map[string]string) *Producer {
	log.Print("Creating producer with config - ", brokers, " and ", topic)
	p, err := kafka.NewProducer(getKafkaConfigMap(brokers, configs))
	if err != nil {
		log.Fatal("Failed to create producer: ", err)
	}
	return &Producer{
		producer: p,
		Topic:    topic,
		Configs:  configs,
		Brokers:  brokers,
	}
}

// Close to close the kafka producer
func (p *Producer) Close() {
	p.producer.Close()
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
	deliveryEvent := <-p.producer.Events()
	return checkMessageDelivery(deliveryEvent)
}

func checkMessageDelivery(deliveryEvent kafka.Event) error {
	switch event := deliveryEvent.(type) {
	case *kafka.Message:
		if event.TopicPartition.Error != nil {
			log.Printf("Delivery failed: %v\n", event.TopicPartition.Error)
			return event.TopicPartition.Error
		}
		log.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*event.TopicPartition.Topic, event.TopicPartition.Partition, event.TopicPartition.Offset)
	case *kafka.Error:
		log.Printf("Delivery failed: %v\n", event.Error())
		return errors.New(event.Error())
	default:
		log.Printf("Ignoring event: %s\n", event)
	}
	return nil
}
