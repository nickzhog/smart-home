package repository

import (
	"context"
	"log/slog"
	"strings"
	"temperature-manager/pkg/postgres"
)

type repository struct {
	client postgres.Client
	logger *slog.Logger
}

func NewRepository(client postgres.Client, logger *slog.Logger) *repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func (r *repository) UpdateCurrentTemperature(ctx context.Context, id string, temp int) error {
	q := `
	UPDATE 
		sensors
	SET 
		current_temperature = $1 
	WHERE 
		id = $2
	`
	r.logger.Info("update current temperature", "query", format(q))

	_, err := r.client.Exec(ctx, q, id, temp)

	if err != nil {
		r.logger.Error("cant update", "error", err, "id", id)
		return err
	}

	return nil
}

func (r *repository) UpdateTargetTemperature(ctx context.Context, id string, temp int) error {
	q := `
	UPDATE 
		sensors
	SET 
		target_temperature = $1 
	WHERE 
		id = $2
	`
	r.logger.Info("update target request", "query", format(q))

	_, err := r.client.Exec(ctx, q, id, temp)

	if err != nil {
		r.logger.Error("cant update", "error", err, "id", id)
		return err
	}

	return nil
}

func (r *repository) FindTargetTemperature(ctx context.Context, sensorId string) (*int, error) {
	q := `
		SELECT
		 	target_temperature
		FROM 
			public.sensors 
		WHERE 
			id = $1;
	`
	r.logger.Info("find target temperature", "query", format(q))
	var target int
	err := r.client.QueryRow(ctx, q, sensorId).Scan(&target)
	if err != nil {
		r.logger.Error("cant find", "error", err, "id", sensorId)
		return nil, err
	}
	return &target, nil
}
func (r *repository) FindCurrentTemperature(ctx context.Context, sensorId string) (*int, error) {
	q := `
		SELECT
		 	current_temperature
		FROM 
			public.sensors 
		WHERE 
			id = $1;
	`
	r.logger.Info("find target temperature", "query", format(q))
	var target int
	err := r.client.QueryRow(ctx, q, sensorId).Scan(&target)
	if err != nil {
		r.logger.Error("current temp find err", "error", err, "id", sensorId)
		return nil, err
	}
	return &target, nil
}

func format(s string) string {
	return strings.ReplaceAll(s, " ", "")
}
