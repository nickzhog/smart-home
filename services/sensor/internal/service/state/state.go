package state

import "sync"

type State struct {
	currentTemperature int
	targetTemperature  int

	mu *sync.RWMutex
}

func NewSensorState() *State {
	return &State{}
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
}
func (s *State) SetCurrentTemperature(val int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentTemperature = val
}
