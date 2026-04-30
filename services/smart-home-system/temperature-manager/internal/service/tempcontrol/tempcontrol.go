package tempcontrol

import (
	"context"
	"log/slog"
)

type TemperatureEventsProducer interface {
	ChangedTargetTemp(ctx context.Context, sensorId string, val int) error
}
type Storage interface {
	FindCurrentTemperature(ctx context.Context, sensorId string) (*int, error)
	FindTargetTemperature(ctx context.Context, sensorId string) (*int, error)
	UpdateCurrentTemperature(ctx context.Context, id string, temp int) error
	UpdateTargetTemperature(ctx context.Context, id string, temp int) error
}
type tempController struct {
	producer TemperatureEventsProducer
	storage  Storage
	logger   *slog.Logger
}

func NewController(logger *slog.Logger, producer TemperatureEventsProducer, storage Storage) *tempController {
	return &tempController{
		producer: producer,
		storage:  storage,
		logger:   logger,
	}
}

func (c *tempController) FindSensorCurrentTemp(ctx context.Context, sensorId string) (*int, error) {
	return c.storage.FindCurrentTemperature(ctx, sensorId)
}
func (c *tempController) FindSensorTargetTemp(ctx context.Context, sensorId string) (*int, error) {
	return c.storage.FindTargetTemperature(ctx, sensorId)
}

func (c *tempController) ChangeSensorCurrentTemp(ctx context.Context, sensorId string, val int) error {
	err := c.storage.UpdateCurrentTemperature(ctx, sensorId, val)
	if err != nil {
		return err
	}
	return nil
}
func (c *tempController) ChangeSensorTargetTemp(ctx context.Context, sensorId string, val int) error {
	err := c.storage.UpdateTargetTemperature(ctx, sensorId, val)
	if err != nil {
		return err
	}
	return c.producer.ChangedTargetTemp(ctx, sensorId, val)
}
