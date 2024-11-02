package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sensor/internal/client"
	"sensor/internal/server"
	"sensor/internal/service/state"
	"sensor/pkg/logging"
	"syscall"
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

	logger = logger.With(slog.String("sensor_id", cfg.SensorId))

	logger.Info(fmt.Sprintf("cfg: %+v", *cfg))

	sensorState := state.NewSensorState()

	srv := server.NewServer(
		logger.With(slog.String("http_server_interface", cfg.WebServerListen)),
		sensorState)
	go srv.Listen(ctx, cfg.WebServerListen)

	client := client.NewClient(
		logger.With(slog.String("client_addr", cfg.SmartHomeApiAddress)),
		sensorState, cfg.SensorId)
	go client.Start(ctx, cfg.SmartHomeApiAddress)

	<-ctx.Done()
	logger.Info("shutdown gracefull")
}
