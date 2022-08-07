package webserver

//go:generate static ./../../web/dist/fdns ui

import (
	"net/http"

	"github.com/deweppro/go-logger"
	"github.com/deweppro/go-static"
)

var ui static.Reader

func (v *WebServer) RegisterUI() {
	for _, file := range ui.List() {
		logger.Debugf("static: %s", file)
		v.route.Route(file, v.Static, http.MethodGet)
	}
	v.route.Route("/", v.Static, http.MethodGet)
}

//Static controller
func (v *WebServer) Static(w http.ResponseWriter, r *http.Request) {
	filename := r.RequestURI
	switch filename {
	case "", "/":
		filename = "/index.html"
		break
	}

	if err := ui.ResponseWrite(w, filename); err != nil {
		logger.Errorf("static response: %s", err.Error())
	}
}
