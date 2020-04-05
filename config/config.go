package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

// Config stores configurations for the application
type Config struct {
	Server struct {
		Port int    `yaml:"port" env:"SERVER_PORT"`
		Host string `yaml:"host" env:"SERVER_HOST"`
	} `yaml:"server"`
}

// Configs initializes configuration from a YAML file and then loads overrides from env
func Configs(configPath string) *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}
