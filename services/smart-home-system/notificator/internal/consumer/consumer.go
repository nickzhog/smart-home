package consumer

import (
	"context"
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
	// config.Version = sarama.V2_5_0_0
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
