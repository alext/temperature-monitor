package webserver

import (
	"encoding/json"
	"fmt"
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
	r.Methods("PUT").Path("/sensors/{sensor_id}").HandlerFunc(srv.sensorPut)
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

func (srv *Webserver) sensorPut(w http.ResponseWriter, req *http.Request) {
	sensorID := mux.Vars(req)["sensor_id"]
	s, ok := srv.sensors[sensorID]
	if !ok {
		write404(w)
		return
	}
	ss, ok := s.(sensor.SettableSensor)
	if !ok {
		writeError(w, fmt.Errorf("Non-writable sensor %s", sensorID), http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Temp *int `json:"temperature"`
	}
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if data.Temp == nil {
		writeError(w, fmt.Errorf("Missing temperature data in request"), http.StatusBadRequest)
		return
	}

	ss.Set(*data.Temp, time.Now())

	writeJSON(w, newJSONSensor(ss))
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
