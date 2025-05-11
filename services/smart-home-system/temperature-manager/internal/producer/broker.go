package producer

import (
	"context"
	"encoding/json"
	"log/slog"
	"temperature-manager/internal/service/tempcontrol"
	"temperature-manager/pkg/logging"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

type producer struct {
	logger    *slog.Logger
	topic     string
	kafkaAddr string

	kafkaProducer *nsq.Producer
}

var _ tempcontrol.TemperatureEventsProducer = (*producer)(nil)

var tries = 5

func New(logger *slog.Logger, brokerAddr, topic string) (*producer, error) {
	config := nsq.NewConfig()
	nsqproducer, err := nsq.NewProducer(brokerAddr, config)
	if err != nil {
		logger.Error("failed to create producer", "error", err)
		if tries != 0 {
			tries--
			time.Sleep(time.Second)
			logger.Error("producer init error, try again", logging.ErrAttr(err))
			return New(logger, brokerAddr, topic)
		}
		return nil, err
	}
	// defer producer.Close()

	producer := &producer{
		logger: logger,

		topic:     topic,
		kafkaAddr: brokerAddr,

		kafkaProducer: nsqproducer,
	}
	return producer, nil
}

type EventTargetTemperature struct {
	SensorId string `json:"sensor_id,omitempty"`
	Value    int    `json:"value,omitempty"`
}

func (p *producer) ChangedTargetTemp(ctx context.Context, sensorId string, val int) error {
	event := EventTargetTemperature{
		SensorId: sensorId,
		Value:    val,
	}
	eventJson, err := json.Marshal(event)
	if err != nil {
		p.logger.Error("json error", logging.ErrAttr(err))
		return err
	}
	err = p.kafkaProducer.Publish(p.topic, eventJson)
	if err != nil {
		p.logger.Error("cant send to kafka", logging.ErrAttr(err))
		return err
	}
	p.logger.Info("sended to kafka target temp change",
		"sensor", sensorId,
		"value", val,
	)
	return nil
}
