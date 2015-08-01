package webserver_test

// dummy sensor implementation for testing webserver
type dummySensor struct {
	temp float64
}

func (s *dummySensor) Temperature() float64 {
	return s.temp
}

func (s *dummySensor) Close() {
}

func (s *dummySensor) SetTemperature(value float64) {
	s.temp = value
}
