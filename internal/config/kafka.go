package config

type Kafka struct {
	Addr         string `yaml:"addr"`
	OrdersTopic  string `env:"orders_topic"`
	ServiceGroup string `env:"service-group"`
}
