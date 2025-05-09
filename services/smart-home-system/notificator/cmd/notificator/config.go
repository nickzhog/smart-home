package main

import "github.com/caarlos0/env"

type config struct {
	KafkaAddr      string `env:"KAFKA_ADDR"`
	KafkaTopic     string `env:"KAFKA_TOPIC"`
	KafkaPartition int    `env:"KAFKA_PARTITION"`
	KafkaOffset    int    `env:"KAFKA_OFFSET"`
}

func parseConfig() (*config, error) {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
