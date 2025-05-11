package main

import (
	"context"
	"fmt"
	"log/slog"
	"notificator/internal/client"
	"notificator/internal/consumer"
	"notificator/pkg/logging"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := parseConfig()
	if err != nil {
		logger.Error("cant parse config", logging.ErrAttr(err))
		os.Exit(255)
	}
	logger.Info(fmt.Sprintf("notificator cfg: %+v", *cfg))

	client := client.NewNotifier(logger.With(slog.Time("http_client_init_time", time.Now())))
	_, err = consumer.New(
		logger.With(slog.Time("consumer_init_time", time.Now())),
		client,
		cfg.KafkaAddr,
		cfg.KafkaTopic,
		cfg.KafkaPartition,
		cfg.KafkaOffset)
	if err != nil {
		logger.Error("cant consumer init", logging.ErrAttr(err))
		os.Exit(255)
	}
	logger.Info("starting")
	<-context.Background().Done()
}
