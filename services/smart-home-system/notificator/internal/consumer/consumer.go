package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"notificator/pkg/logging"
	"os"
	"time"

	"github.com/IBM/sarama"
)

type Sender interface {
	SendTargetTemperatureChangeEvent(ctx context.Context, sensorId string, val int) error
}
type consumer struct {
	logger       *slog.Logger
	partConsumer sarama.PartitionConsumer
	sender       Sender
}

var tries = 5000000

func New(logger *slog.Logger, sender Sender, kafkaAddr, topic string) (*consumer, error) {
	consumerKafka, err := sarama.NewConsumer([]string{kafkaAddr}, nil)
	if err != nil {
		if tries != 0 {
			time.Sleep(time.Second)
			tries--
			logger.Error("consumer init error, try again", logging.ErrAttr(err))
			return New(logger, sender, kafkaAddr, topic)
		}
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}
	// defer consumer.Close()
	partConsumer, err := consumerKafka.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		if tries != 0 {
			time.Sleep(time.Second)
			tries--
			logger.Error("consumer(partition) init error, try again", logging.ErrAttr(err))
			return New(logger, sender, kafkaAddr, topic)
		}
		return nil, fmt.Errorf("failed to consume partition: %w", err)
	}
	// defer partConsumer.Close()
	consumer := &consumer{
		logger:       logger,
		partConsumer: partConsumer,
		sender:       sender,
	}
	return consumer, nil
}

type EventTargetTemperature struct {
	SensorId string `json:"sensor_id,omitempty"`
	Value    int    `json:"value,omitempty"`
}

func (c *consumer) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("stop")
		case msg, ok := <-c.partConsumer.Messages():
			if !ok {
				c.logger.Error("partition channel closed")
				os.Exit(255)
			}
			var event EventTargetTemperature
			err := json.Unmarshal(msg.Value, &event)
			if err != nil {
				c.logger.Error("cant parse event",
					logging.ErrAttr(err),
					"data", string(msg.Value))
				continue
			}
			err = c.sender.SendTargetTemperatureChangeEvent(ctx, event.SensorId, event.Value)
			if err != nil {
				c.logger.Error("cant handle event",
					logging.ErrAttr(err),
					"data", string(msg.Value))
				continue
			}
			c.logger.Info("handled event",
				"event", fmt.Sprintf("%+v", event))
		}
	}
}
