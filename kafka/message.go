package kafka

import (
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Message represents each kafka message and it's metadata
type Message struct {
	kafkaMessage *kafka.Message
}

// GetMessage will get the kafka message as a byte array
func (m *Message) GetMessage() *[]byte {
	return &m.kafkaMessage.Value
}

// GetKey will get the kafka key as a byte array
func (m *Message) GetKey() *[]byte {
	return &m.kafkaMessage.Key
}

// GetTopic will get the kafka topic string ref
func (m *Message) GetTopic() *string {
	return m.kafkaMessage.TopicPartition.Topic
}

// GetPartition will get the kafka message partition
func (m *Message) GetPartition() int32 {
	return m.kafkaMessage.TopicPartition.Partition
}

// GetOffset will get the kafka message offset
func (m *Message) GetOffset() int64 {
	return int64(m.kafkaMessage.TopicPartition.Offset)
}

// GetTimestamp will get the kafka message timestamp
func (m *Message) GetTimestamp() *time.Time {
	return &m.kafkaMessage.Timestamp
}
