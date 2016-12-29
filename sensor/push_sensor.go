package sensor

import (
	"sync"
	"time"
)

type pushSensor struct {
	sensorID string

	mu        sync.RWMutex
	temp      int
	updatedAt time.Time
}

func NewPushSensor(id string) SettableSensor {
	return &pushSensor{
		sensorID: id,
	}
}

func (s *pushSensor) Read() (int, time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.temp, s.updatedAt
}

func (s *pushSensor) Close() {
	// No-Op
}

func (s *pushSensor) Set(temp int, updatedAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.temp = temp
	s.updatedAt = updatedAt
}
