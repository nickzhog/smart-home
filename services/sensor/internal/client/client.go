// send interval requests with current temperature if changed
package client

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sensor/internal/service/state"
	"sensor/pkg/logging"
	"strings"
	"time"
)

type Client struct {
	logger             *slog.Logger
	currentTemperature int
	state              *state.State
}

func NewClient(logger *slog.Logger, state *state.State) *Client {
	c := &Client{}

	return c
}

func (c *Client) Start(ctx context.Context, addr string) {
	addr = strings.TrimSuffix(addr, "/")

	c.currentTemperature = c.state.GetCurrentTemperature()

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			actualCur := c.state.GetCurrentTemperature()
			if actualCur == c.currentTemperature {
				continue
			}
			beforeCur := c.currentTemperature
			c.currentTemperature = actualCur

			url := addr + "/current_temperature/" + fmt.Sprint(actualCur)
			c.logger.Info("current temperature changes",
				slog.Int("actual", actualCur),
				slog.Int("before", beforeCur),
				slog.String("url", url))

			resp, err := http.Post(url, "application/text", nil)
			if err != nil {
				c.logger.Error("with webhook",
					logging.ErrAttr(err),
					slog.Int("actual", actualCur),
					slog.Int("before", beforeCur),
					slog.String("url", url))
				continue
			}
			c.logger.Info("response code", slog.Int("code", resp.StatusCode))

		case <-ctx.Done():
			c.logger.Info("shutdown")
			return
		}
	}
}
