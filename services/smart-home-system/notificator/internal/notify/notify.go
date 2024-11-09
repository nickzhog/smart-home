package notify

import (
	"fmt"
	"log/slog"
	"net/http"
	"notificator/pkg/logging"
	"strings"
)

type Notifier struct {
	webhookUrl string
	logger     *slog.Logger
}

func NewNotifier(logger *slog.Logger, webhookAddr string) *Notifier {
	return &Notifier{
		webhookUrl: webhookAddr,
		logger:     logger,
	}
}

func (n *Notifier) SendTargetTemperatureToSensor(sensorAddr string, temperature int) error {
	sensorAddr = strings.TrimSuffix(sensorAddr, "/")
	url := sensorAddr + "/target/" + fmt.Sprint(temperature)
	response, err := http.Post(url, "application/text", nil)
	if err != nil {
		n.logger.Error("cant send target temp to sensor", logging.ErrAttr(err), "url", url)
		return err
	}
	defer response.Body.Close()
	n.logger.Info("sended target temp to sensor", "url", url, "response_code", response.StatusCode)
	return nil
}
