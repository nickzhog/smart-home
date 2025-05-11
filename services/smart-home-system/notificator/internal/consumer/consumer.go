package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"notificator/pkg/logging"
	"time"

	"github.com/nsqio/go-nsq"
)

type Sender interface {
	SendTargetTemperatureChangeEvent(ctx context.Context, sensorId string, val int) error
}
type consumer struct {
	logger       *slog.Logger
	partConsumer *nsq.Consumer
	sender       Sender

	topic string
}

var tries = 100

func New(logger *slog.Logger, sender Sender, kafkaAddr, topic string, partition, offset int) (*consumer, error) {
	config := nsq.NewConfig()
	nsqconsumer, err := nsq.NewConsumer(topic, fmt.Sprintf("test-channel-%v", partition), config)
	if err != nil {
		if tries != 0 {
			time.Sleep(time.Second)
			tries--
			logger.Error("consumer init error, try again", logging.ErrAttr(err))
			return New(logger, sender, kafkaAddr, topic, partition, offset)
		}
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	consumer := &consumer{
		logger:       logger,
		partConsumer: nsqconsumer,
		sender:       sender,
		topic:        topic,
	}

	nsqconsumer.AddHandler(&handler{c: *consumer})
	// Use nsqlookupd to discover nsqd instances
	err = nsqconsumer.ConnectToNSQLookupd(kafkaAddr)
	if err != nil {
		if tries != 0 {
			time.Sleep(time.Second)
			tries--
			logger.Error("consumer init error, try again", logging.ErrAttr(err))
			return New(logger, sender, kafkaAddr, topic, partition, offset)
		}
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return consumer, nil
}

type EventTargetTemperature struct {
	SensorId string `json:"sensor_id,omitempty"`
	Value    int    `json:"value,omitempty"`
}

type handler struct {
	c consumer
}

func (h *handler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		h.c.logger.Error("got msg with empty body")
		return nil
	}

	var event EventTargetTemperature
	err := json.Unmarshal(m.Body, &event)
	if err != nil {
		h.c.logger.Error("cant parse event",
			logging.ErrAttr(err),
			"data", string(m.Body))
		return err
	}

	// эмуляция нагрузки
	time.Sleep(time.Second * 2)
	//

	err = h.c.sender.SendTargetTemperatureChangeEvent(context.TODO(), event.SensorId, event.Value)
	if err != nil {
		h.c.logger.Error("cant handle event",
			logging.ErrAttr(err),
			"data", string(m.Body))
		return err
	}

	h.c.logger.Info("handled event",
		"event", fmt.Sprintf("%+v", event))
	m.Finish()

	return nil
}
