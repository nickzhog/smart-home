package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"notificator/pkg/logging"
	"time"

	"github.com/IBM/sarama"
)

type Sender interface {
	SendTargetTemperatureChangeEvent(ctx context.Context, sensorId string, val int) error
}
type consumer struct {
	logger       *slog.Logger
	partConsumer sarama.ConsumerGroup
	sender       Sender

	topic string
}

var tries = 5000000

func New(logger *slog.Logger, sender Sender, kafkaAddr, topic string, partition, offset int) (*consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup([]string{kafkaAddr}, "groupID", config)
	if err != nil {
		if tries != 0 {
			time.Sleep(time.Second)
			tries--
			logger.Error("consumer init error, try again", logging.ErrAttr(err))
			return New(logger, sender, kafkaAddr, topic, partition, offset)
		}
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}
	// defer group.Close()

	consumer := &consumer{
		logger:       logger,
		partConsumer: group,
		sender:       sender,
		topic:        topic,
	}
	return consumer, nil
}

type EventTargetTemperature struct {
	SensorId string `json:"sensor_id,omitempty"`
	Value    int    `json:"value,omitempty"`
}

func (c *consumer) Listen(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			c.logger.Info("exit by ctx")
			return
		}
		err := c.partConsumer.Consume(ctx, []string{c.topic}, &ConsumerHandler{ctx: ctx, c: c})
		if err != nil {
			c.logger.Error("consume err", logging.ErrAttr(err))
			continue
		}
	}
}

type ConsumerHandler struct {
	ctx context.Context
	c   *consumer
}

func (h *ConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	h.c.logger.Info("New session started",
		"claims", session.Claims(),
		"member_id", session.MemberID())
	return nil

}
func (h *ConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	h.c.logger.Info("cleanup",
		"claims", session.Claims(),
		"member_id", session.MemberID())
	return nil
}
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	h.c.logger.Info("Starting consumption for partition",
		"partition", claim.Partition(),
		"initial_offset", claim.InitialOffset())

	defer func() {
		h.c.logger.Info("Stopping consumption for partition",
			"partition", claim.Partition())
	}()

	for {
		select {
		case <-h.ctx.Done():
			h.c.logger.Info("stop")
			return nil

		case msg, ok := <-claim.Messages():
			if !ok {
				h.c.logger.Error("partition channel closed")
				return fmt.Errorf("partition channel closed")
			}
			var event EventTargetTemperature
			err := json.Unmarshal(msg.Value, &event)
			if err != nil {
				h.c.logger.Error("cant parse event",
					logging.ErrAttr(err),
					"data", string(msg.Value))
				continue
			}

			// эмуляция нагрузки
			time.Sleep(time.Second * 2)
			//

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
