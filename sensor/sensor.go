package sensor

import "github.com/spf13/afero"

var fs afero.Fs = &afero.OsFs{}

type Sensor interface {
	Temperature() float64
}
