package webserver

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (srv *Webserver) buildHandler() http.Handler {
	r := mux.NewRouter()
	r.Methods("GET").Path("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("OK\n"))
	})
	r.Methods("GET").Path("/sensors/{sensor_id}").HandlerFunc(srv.sensorGet)
	return r
}

func (srv *Webserver) sensorGet(w http.ResponseWriter, req *http.Request) {
	s, ok := srv.sensors[mux.Vars(req)["sensor_id"]]
	if !ok {
		write404(w)
		return
	}

	jsonData, err := json.MarshalIndent(map[string]interface{}{
		"temperature": s.Temperature(),
	}, "", "  ")
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
