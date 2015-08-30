package sensor

import (
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/alext/temperature-monitor/Godeps/_workspace/src/github.com/alext/afero"
)

var fs afero.Fs = &afero.OsFs{}

const w1DevicesPath = "/sys/bus/w1/devices/"

type Sensor interface {
	Read() (int, time.Time)
	Close()
}

type sensor struct {
	deviceID  string
	mux       sync.RWMutex
	temp      int
	updatedAt time.Time
	closeCh   chan struct{}
}

func New(deviceID string) (Sensor, error) {
	s := &sensor{
		deviceID: deviceID,
		closeCh:  make(chan struct{}),
	}
	s.readTemperature(time.Now())
	go s.readLoop()
	return s, nil
}

func (s *sensor) readLoop() {
	t := newTicker(time.Minute)
	for {
		select {
		case t := <-t.Channel():
			s.readTemperature(t)
		case <-s.closeCh:
			t.Stop()
			close(s.closeCh)
			return
		}
	}
}

func (s *sensor) Read() (int, time.Time) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.temp, s.updatedAt
}

func (s *sensor) Close() {
	s.closeCh <- struct{}{}
	<-s.closeCh
}

var temperatureRegexp = regexp.MustCompile(`t=(\d+)`)

func (s *sensor) readTemperature(updateTime time.Time) {
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

	temp, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Printf("[sensor:%s] Error parsing temperature value '%s': %s", s.deviceID, matches[1], err.Error())
		return
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	s.temp = temp
	s.updatedAt = updateTime
}
