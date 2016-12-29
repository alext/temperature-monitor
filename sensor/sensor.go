package sensor

import (
	"time"

	"github.com/spf13/afero"
)

var fs afero.Fs = &afero.OsFs{}

type Sensor interface {
	Read() (int, time.Time)
	Close()
}

type SettableSensor interface {
	Sensor
	Set(int, time.Time)
}
