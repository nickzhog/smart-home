package producer

import (
	"context"
	"encoding/json"
	"log/slog"
	"temperature-manager/internal/service/tempcontrol"
	"temperature-manager/pkg/logging"
	"time"

	"github.com/IBM/sarama"
)

type producer struct {
	logger    *slog.Logger
	topic     string
	kafkaAddr string

	kafkaProducer sarama.SyncProducer
}

var _ tempcontrol.TemperatureEventsProducer = (*producer)(nil)

func New(logger *slog.Logger, kafkaAddr, topic string) (*producer, error) {
	time.Sleep(time.Second * 5)
	kafkaProducer, err := sarama.NewSyncProducer([]string{kafkaAddr}, nil)
	if err != nil {
		logger.Error("Failed to create producer", "error", err)
		return nil, err
	}
	// defer producer.Close()

	producer := &producer{
		logger: logger,

		topic:     topic,
		kafkaAddr: kafkaAddr,

		kafkaProducer: kafkaProducer,
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
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder("target_temperature_changes"),
		Value: sarama.ByteEncoder(eventJson),
	}
	_, _, err = p.kafkaProducer.SendMessage(msg)
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
