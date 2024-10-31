package main

import "github.com/caarlos0/env"

type config struct {
	SmartHomeApiAddress string `env:"SMART_HOME_ADDRESS"`
	WebServerListen     string `env:"WEBSERVER_LISTEN"`
	SensorId            string `env:"SENSOR_ID"`
}

func parseConfig() (*config, error) {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
