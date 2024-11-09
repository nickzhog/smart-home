package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"temperature-manager/pkg/logging"

	"github.com/go-chi/chi"
)

func (s *server) changeTargetTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	sensorId := chi.URLParam(r, "sensor_id")
	if len(sensorId) == 0 {
		http.Error(w, "empty sensor", http.StatusBadRequest)
		return
	}
	value := chi.URLParam(r, "value")
	if len(value) == 0 {
		http.Error(w, "empty value", http.StatusBadRequest)
		return
	}
	s.logger.Info("change target", slog.String("sensor", sensorId), slog.String("value", value))
	targetTemp, err := strconv.Atoi(value)
	if err != nil {
		s.logger.Error("cant parse value", logging.ErrAttr(err))
		http.Error(w, "cant parse value: "+err.Error(), http.StatusBadRequest)
		return
	}
	err = s.tempController.ChangeSensorTargetTemp(r.Context(), sensorId, targetTemp)
	if err != nil {
		s.logger.Error("cant change temp", logging.ErrAttr(err))
		http.Error(w, "cant change temp: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "ok")
}

func (s *server) changeCurrentTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	sensorId := chi.URLParam(r, "sensor_id")
	if len(sensorId) == 0 {
		http.Error(w, "empty sensor", http.StatusBadRequest)
		return
	}
	value := chi.URLParam(r, "value")
	if len(value) == 0 {
		http.Error(w, "empty value", http.StatusBadRequest)
		return
	}
	s.logger.Info("change current", slog.String("sensor", sensorId), slog.String("value", value))

	curTemp, err := strconv.Atoi(value)
	if err != nil {
		s.logger.Error("cant parse value", logging.ErrAttr(err))
		http.Error(w, "cant parse value: "+err.Error(), http.StatusBadRequest)
		return
	}
	err = s.tempController.ChangeSensorCurrentTemp(r.Context(), sensorId, curTemp)
	if err != nil {
		s.logger.Error("cant change temp", logging.ErrAttr(err))
		http.Error(w, "cant change temp: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "ok")
}

func (s *server) getCurrentTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	sensorId := chi.URLParam(r, "sensor_id")
	if len(sensorId) == 0 {
		http.Error(w, "empty sensor", http.StatusBadRequest)
		return
	}
	s.logger.Info("get current", slog.String("sensor", sensorId))

	val, err := s.tempController.FindSensorCurrentTemp(r.Context(), sensorId)
	if err != nil {
		s.logger.Error("cant get cur temp", logging.ErrAttr(err))
		http.Error(w, "cant get cur temp: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, *val)
}
func (s *server) getTargetTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	sensorId := chi.URLParam(r, "sensor_id")
	if len(sensorId) == 0 {
		http.Error(w, "empty sensor", http.StatusBadRequest)
		return
	}
	s.logger.Info("get target", slog.String("sensor", sensorId))

	val, err := s.tempController.FindSensorTargetTemp(r.Context(), sensorId)
	if err != nil {
		s.logger.Error("cant get tar temp", logging.ErrAttr(err))
		http.Error(w, "cant get tar temp: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, *val)
}
