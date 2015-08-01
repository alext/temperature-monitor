package webserver

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (srv *Webserver) buildHandler() http.Handler {
	r := mux.NewRouter()
	r.Methods("GET").Path("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("OK\n"))
	})
	return r
}
