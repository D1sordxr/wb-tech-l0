package config

import "time"

type Kafka struct {
	Address          string        `yaml:"address" env:"KAFKA_ADDRESS"`
	OrdersTopic      string        `yaml:"orders_topic" env:"KAFKA_ORDERS_TOPIC"`
	SaverGroup       string        `yaml:"saver_group" env:"KAFKA_SAVER_GROUP"`
	BroadcasterGroup string        `yaml:"broadcaster_group" env:"KAFKA_BROADCASTER_GROUP"`
	CreateTopic      bool          `yaml:"create_topic" env:"KAFKA_CREATE_TOPIC"`
	SessionTimeout   time.Duration `yaml:"session_timeout" env:"KAFKA_SESSION_TIMEOUT" env-default:"30s"`
	MaxPollInterval  time.Duration `yaml:"max_poll_interval" env:"KAFKA_MAX_POLL_INTERVAL" env-default:"5m"`
}
