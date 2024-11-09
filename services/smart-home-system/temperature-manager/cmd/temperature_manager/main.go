package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"temperature-manager/internal/producer"
	"temperature-manager/internal/repository"
	"temperature-manager/internal/server"
	"temperature-manager/internal/service/tempcontrol"
	"temperature-manager/pkg/logging"
	"temperature-manager/pkg/postgres"
	"time"
)

func main() {
	exitSignalCh := make(chan os.Signal, 1)
	signal.Notify(exitSignalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-exitSignalCh
		cancel()
	}()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := parseConfig()
	if err != nil {
		logger.Error("cant parse config", logging.ErrAttr(err))
		os.Exit(255)
	}

	logger.Info(fmt.Sprintf("manager cfg: %+v", *cfg))

	postgresClient, err := postgres.NewClient(ctx, 5, cfg.DatabaseUri)
	if err != nil {
		logger.Error("cant connect to db", "error", err)
		panic(err)
	}
	rep := repository.NewRepository(
		postgresClient,
		logger.With(slog.Time("db_init_time", time.Now())),
	)
	producer, err := producer.New(
		logger.With(slog.Time("producer_init_time", time.Now())),
		cfg.KafkaAddr,
		cfg.KafkaTopic,
	)
	if err != nil {
		logger.Error("cant connect to producer", "error", err)
		panic(err)
	}

	tempController := tempcontrol.NewController(
		logger.With(slog.Time("controller_init_time", time.Now())),
		producer,
		rep,
	)

	srv := server.NewHttpServer(
		logger.With(slog.Time("http_server_init_time", time.Now())),
		tempController,
	)

	wg := new(sync.WaitGroup)
	go func() {
		wg.Add(1)
		srv.Serve(ctx)
		wg.Done()
	}()
	wg.Wait()
	logger.Info("program stopped")
}
