package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type TempController interface {
	FindSensorCurrentTemp(ctx context.Context, sensorId string) (*int, error)
	FindSensorTargetTemp(ctx context.Context, sensorId string) (*int, error)

	ChangeSensorCurrentTemp(ctx context.Context, sensorId string, val int) error
	ChangeSensorTargetTemp(ctx context.Context, sensorId string, val int) error
}
type server struct {
	logger         *slog.Logger
	tempController TempController
	srv            *http.Server
}

func NewHttpServer(logger *slog.Logger, controller TempController, addr string) *server {
	server := &server{
		logger:         logger,
		tempController: controller,
	}
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	logger.Info("asdasdasd")

	r.Route("/", func(r chi.Router) {
		r.Route("/sensor/{sensor_id}", func(r chi.Router) {
			r.Get("/current-temperature", server.getCurrentTemperatureHandler)
			r.Get("/target-temperature", server.getTargetTemperatureHandler)
			r.Post("/target-temperature/{value}", server.changeTargetTemperatureHandler)
			r.Post("/current-temperature/{value}", server.changeCurrentTemperatureHandler)
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			fmt.Fprintln(w, "not found")
		})
	})

	server.srv = &http.Server{
		Addr:     addr,
		Handler:  r,
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	return server
}

func (s *server) Serve(ctx context.Context) {
	go func() {
		for {
			if ctx.Err() != nil {
				s.logger.Info("context end")
			}
			if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				s.logger.Error("listen error", "error", err)
			}
		}
	}()

	s.logger.Info("server started")

	<-ctx.Done()

	s.logger.Info("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctxShutDown); err != nil {
		s.logger.Error("server Shutdown Failed", "error", err)
	}

	s.logger.Info("server exited properly")
}
