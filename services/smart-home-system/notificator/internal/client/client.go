package client

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"notificator/pkg/logging"
)

type Notifier struct {
	logger *slog.Logger
}

func NewNotifier(logger *slog.Logger) *Notifier {
	return &Notifier{
		logger: logger,
	}
}
func (n *Notifier) SendTargetTemperatureChangeEvent(ctx context.Context, sensorId string, val int) error {
	url := fmt.Sprintf("http://%s/target/%v", sensorId, val)
	response, err := http.Post(url, "application/text", nil)
	if err != nil {
		n.logger.Error("cant send target temp to sensor", logging.ErrAttr(err), "url", url)
		return err
	}
	defer response.Body.Close()
	n.logger.Info("sended target temp to sensor",
		"url", url,
		"response_code", response.StatusCode)
	return nil
}
