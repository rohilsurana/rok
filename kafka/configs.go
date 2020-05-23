package kafka

import "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

func getKafkaConfigMap(bootstrapServers string, configs map[string]string) *kafka.ConfigMap {
	configs["bootstrap.servers"] = bootstrapServers
	kafkaConfigs := kafka.ConfigMap{}
	for k, v := range configs {
		kafkaConfigs[k] = v
	}
	return &kafkaConfigs
}
