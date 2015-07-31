package sensor

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"

	"github.com/spf13/afero"
)

var fs afero.Fs = &afero.OsFs{}

const w1DevicesPath = "/sys/bus/w1/devices/"

type Sensor interface {
	Temperature() float64
	Close()
}

type sensor struct {
	deviceID string
	temp     float64
	closeCh  chan struct{}
}

func New(deviceID string) (Sensor, error) {
	s := &sensor{
		deviceID: deviceID,
		closeCh:  make(chan struct{}),
	}
	s.readTemperature()
	go s.readLoop()
	return s, nil
}

func (s *sensor) readLoop() {
	t := newTicker(time.Minute)
	for {
		select {
		case <-t.Channel():
			s.readTemperature()
		case <-s.closeCh:
			t.Stop()
			close(s.closeCh)
			return
		}
	}
}

func (s *sensor) Temperature() float64 {
	return s.temp
}

func (s *sensor) Close() {
	s.closeCh <- struct{}{}
	<-s.closeCh
}

var temperatureRegexp = regexp.MustCompile(`t=(\d+)`)

func (s *sensor) readTemperature() {
	file, err := fs.Open(w1DevicesPath + s.deviceID + "/w1_slave")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	matches := temperatureRegexp.FindStringSubmatch(string(data))
	if matches == nil {
		fmt.Println("No Match")
		return
	}
	value, err := strconv.Atoi(matches[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	s.temp = float64(value) / 1000
}
