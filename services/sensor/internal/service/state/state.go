package state

import (
	"sync"
	"time"
)

type State struct {
	currentTemperature int
	targetTemperature  int

	mu *sync.RWMutex
}

func NewSensorState() *State {
	return &State{
		mu: new(sync.RWMutex),
	}
}

func (s *State) GetCurrentTemperature() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentTemperature
}

func (s *State) GetTargetTemperature() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.targetTemperature
}

func (s *State) SetTargetTemperature(val int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.targetTemperature = val

	// current temperature change stub
	go func() {
		time.Sleep(time.Second * 5)
		s.SetCurrentTemperature(val)
	}()
}
func (s *State) SetCurrentTemperature(val int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentTemperature = val
}
