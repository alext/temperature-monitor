package webserver_test

// dummy sensor implementation for testing webserver
type dummySensor struct {
	temp int
}

func (s *dummySensor) Temperature() int {
	return s.temp
}

func (s *dummySensor) Close() {
}

func (s *dummySensor) SetTemperature(value int) {
	s.temp = value
}
