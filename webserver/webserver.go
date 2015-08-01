package webserver

import (
	"fmt"
	"net/http"
)

type Webserver struct {
	listenURL string
	mux       http.Handler
}

func New(port int) *Webserver {
	srv := &Webserver{
		listenURL: fmt.Sprintf(":%d", port),
	}
	srv.mux = srv.buildHandler()
	return srv
}

func (srv *Webserver) Run() error {
	return http.ListenAndServe(srv.listenURL, srv)
}

func (srv *Webserver) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	srv.mux.ServeHTTP(w, req)
}
