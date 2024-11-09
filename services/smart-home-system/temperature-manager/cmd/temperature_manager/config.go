package main

import "github.com/caarlos0/env"

type config struct {
	DatabaseUri string `env:"DATABASE_URI"`
	KafkaAddr   string `env:"KAFKA_ADDR"`
	KafkaTopic  string `env:"KAFKA_TOPIC"`
}

func parseConfig() (*config, error) {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
