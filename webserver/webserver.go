package webserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alext/temperature-monitor/sensor"
)

type Webserver struct {
	listenURL string
	mux       http.Handler
	sensors   map[string]sensor.Sensor
}

func New(port int) *Webserver {
	srv := &Webserver{
		listenURL: fmt.Sprintf(":%d", port),
		sensors:   make(map[string]sensor.Sensor),
	}
	srv.mux = srv.buildHandler()
	return srv
}

func (srv *Webserver) AddSensor(name string, s sensor.Sensor) {
	srv.sensors[name] = s
}

func (srv *Webserver) Run() error {
	log.Printf("[webserver] listening on %s", srv.listenURL)
	return http.ListenAndServe(srv.listenURL, srv)
}

func (srv *Webserver) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	srv.mux.ServeHTTP(w, req)
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func write404(w http.ResponseWriter) {
	http.Error(w, "Not found", http.StatusNotFound)
}
func writeError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
