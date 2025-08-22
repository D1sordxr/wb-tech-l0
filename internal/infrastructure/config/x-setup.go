package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const basicConfigPath = "./configs/api/prod.yaml"

type Config struct {
	Server        HttpServer `yaml:"server"`
	MessageBroker Kafka      `yaml:"message_broker"`
	Storage       Postgres   `yaml:"storage"`
}

func NewConfig() *Config {
	var cfg Config

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = basicConfigPath
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}
