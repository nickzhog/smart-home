// server for accept webhooks on target temperature changes
package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"sensor/internal/service/state"
	"sensor/pkg/logging"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

type server struct {
	logger *slog.Logger

	state *state.State

	srv *http.Server
}

func NewServer(logger *slog.Logger, state *state.State) *server {
	server := &server{
		logger: logger,
		state:  state,
	}

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/target/{value}", server.changeTargetTemperatureHandler)
	})

	server.srv = &http.Server{
		Handler:  r,
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	return server
}
func (s *server) Listen(ctx context.Context, listen string) {
	s.logger.Info("start listen")
	s.srv.Addr = listen

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	<-ctx.Done()
	s.logger.Info("server stop")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s.srv.Shutdown(ctxShutDown)

	s.logger.Info("server exited properly")
}

func (s *server) changeTargetTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	value := chi.URLParam(r, "value")
	if len(value) == 0 {
		http.Error(w, "empty value", http.StatusBadRequest)
		return
	}
	s.logger.Info("new target", slog.String("value", value))

	targetTemp, err := strconv.Atoi(value)
	if err != nil {
		s.logger.Error("cant parse value", logging.ErrAttr(err))
		http.Error(w, "cant parse value: "+err.Error(), http.StatusBadRequest)
		return
	}
	s.state.SetTargetTemperature(targetTemp)
}
