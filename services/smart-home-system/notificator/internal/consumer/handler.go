package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"notificator/pkg/logging"
	"os"

	"github.com/IBM/sarama"
)

type ConsumerHandler struct {
	ctx context.Context
	c   *consumer
}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			h.c.logger.Info("stop by session")
			return nil
		case <-h.ctx.Done():
			h.c.logger.Info("stop")
			return nil

		case msg, ok := <-claim.Messages():
			if session.Context().Err() != nil {
				h.c.logger.Info("stop by context")
				return nil
			}
			if !ok {
				h.c.logger.Error("partition channel closed")
				os.Exit(255)
			}
			var event EventTargetTemperature
			err := json.Unmarshal(msg.Value, &event)
			if err != nil {
				h.c.logger.Error("cant parse event",
					logging.ErrAttr(err),
					"data", string(msg.Value))
				continue
			}
			err = h.c.sender.SendTargetTemperatureChangeEvent(h.ctx, event.SensorId, event.Value)
			if err != nil {
				h.c.logger.Error("cant handle event",
					logging.ErrAttr(err),
					"data", string(msg.Value))
				continue
			}
			session.MarkMessage(msg, "") // Подтверждаем обработку

			h.c.logger.Info("handled event",
				"event", fmt.Sprintf("%+v", event))
		}
	}
}
