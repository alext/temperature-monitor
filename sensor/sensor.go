package sensor

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"

	"github.com/spf13/afero"
)

var fs afero.Fs = &afero.OsFs{}

const w1DevicesPath = "/sys/bus/w1/devices/"

type Sensor interface {
	Temperature() float64
}

type sensor struct {
	valueFile afero.File
}

func New(deviceID string) (Sensor, error) {
	file, err := fs.Open(w1DevicesPath + deviceID + "/w1_slave")
	if err != nil {
		return nil, err
	}
	return &sensor{valueFile: file}, nil
}

var temperatureRegexp = regexp.MustCompile(`t=(\d+)`)

func (s *sensor) Temperature() float64 {
	s.valueFile.Seek(0, 0)
	data, err := ioutil.ReadAll(s.valueFile)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	matches := temperatureRegexp.FindStringSubmatch(string(data))
	if matches == nil {
		fmt.Println("No Match")
		return 0
	}
	value, err := strconv.Atoi(matches[1])
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return float64(value) / 1000
}
