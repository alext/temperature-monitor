package webserver

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/alext/temperature-monitor/sensor"
)

func (srv *Webserver) buildHandler() http.Handler {
	r := mux.NewRouter()
	r.Methods("GET").Path("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("OK\n"))
	})
	r.Methods("GET").Path("/sensors").HandlerFunc(srv.sensorIndex)
	r.Methods("GET").Path("/sensors/{sensor_id}").HandlerFunc(srv.sensorGet)
	return r
}

func (srv *Webserver) sensorIndex(w http.ResponseWriter, req *http.Request) {
	data := make(map[string]*jsonSensor)
	for name, s := range srv.sensors {
		data[name] = newJSONSensor(s)
	}
	writeJSON(w, data)
}

func (srv *Webserver) sensorGet(w http.ResponseWriter, req *http.Request) {
	s, ok := srv.sensors[mux.Vars(req)["sensor_id"]]
	if !ok {
		write404(w)
		return
	}

	writeJSON(w, newJSONSensor(s))
}

type jsonSensor struct {
	Temperature int       `json:"temperature"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func newJSONSensor(s sensor.Sensor) *jsonSensor {
	temperature, updatedAt := s.Read()
	return &jsonSensor{
		Temperature: temperature,
		UpdatedAt:   updatedAt,
	}
}
