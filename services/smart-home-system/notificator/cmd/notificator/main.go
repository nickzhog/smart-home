package main

import (
	"fmt"
	"log/slog"
	"notificator/pkg/logging"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := parseConfig()
	if err != nil {
		logger.Error("cant parse config", logging.ErrAttr(err))
		os.Exit(255)
	}

	logger.Info(fmt.Sprintf("notificator cfg: %+v", *cfg))

}
