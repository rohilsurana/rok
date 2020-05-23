package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config stores configurations for the application
type Config struct {
	Mode   string `yaml:"mode" env:"MODE"`
	Server struct {
		Port int    `yaml:"port" env:"SERVER_PORT"`
		Host string `yaml:"host" env:"SERVER_HOST"`
	} `yaml:"server"`
	Kafka struct {
		Topic    string            `yaml:"topic" env:"KAFKA_TOPIC"`
		Brokers  string            `yaml:"brokers" env:"KAFKA_BROKER"`
		Producer map[string]string `yaml:"producer" env:"KAFKA_PRODUCER"`
		Consumer map[string]string `yaml:"Consumer" env:"KAFKA_CONSUMER"`
	} `yaml:"kafka"`
}

// NewConfigs initializes configuration from a YAML file and then loads overrides from env
func NewConfigs(configPath string) *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("Failed to get configs: ", err)
	}
	return &cfg
}
