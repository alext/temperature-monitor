package webserver_test

import "time"

// dummy sensor implementation for testing webserver
type dummySensor struct {
	temp       int
	updateTime time.Time
}

func (s *dummySensor) Read() (int, time.Time) {
	return s.temp, s.updateTime
}

func (s *dummySensor) Close() {
}

func (s *dummySensor) SetTemperature(value int, updated time.Time) {
	s.temp = value
	s.updateTime = updated
}
