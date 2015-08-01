package sensor

import (
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/spf13/afero"
)

var fs afero.Fs = &afero.OsFs{}

const w1DevicesPath = "/sys/bus/w1/devices/"

type Sensor interface {
	Temperature() int
	Close()
}

type sensor struct {
	deviceID string
	temp     int
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

func (s *sensor) Temperature() int {
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
		log.Printf("[sensor:%s] Error opening device file: %s", s.deviceID, err.Error())
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("[sensor:%s] Error reading device: %s", s.deviceID, err.Error())
		return
	}
	matches := temperatureRegexp.FindStringSubmatch(string(data))
	if matches == nil {
		log.Printf("[sensor:%s] Failed to match temperature in data:\n%s", s.deviceID, string(data))
		return
	}

	// discard error because it can't fail due to \d in regexp
	s.temp, _ = strconv.Atoi(matches[1])
}
