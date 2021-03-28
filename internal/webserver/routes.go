package webserver

import (
	"net/http"
	"strings"

	"github.com/deweppro/go-http/web/routes"
	"github.com/deweppro/go-logger"
	"github.com/deweppro/go-static"
)

//go:generate static ./../../web/dist/fdns UI

//UI static archive
var UI = "H4sIAAAAAAAA/2IYBaNgFIxYAAgAAP//Lq+17wAEAAA="

//Routes model
type Routes struct {
	cache *static.Cache
	route *routes.Router
	conf  *MiddlewareConfig
}

//NewRoutes init router
func NewRoutes(conf *MiddlewareConfig) (*Routes, *routes.Router) {
	route := routes.NewRouter()
	return &Routes{
		cache: static.New(),
		route: route,
		conf:  conf,
	}, route
}

//Up startup api service
func (v *Routes) Up() error {
	v.route.Global(routes.RecoveryMiddleware(logger.Default()))
	v.route.Global(routes.ThrottlingMiddleware(v.conf.Middleware.Throttling))

	if err := v.cache.FromBase64TarGZ(UI); err != nil {
		return err
	}

	for _, file := range v.cache.List() {
		logger.Debugf("static: %s", file)
		v.route.Route(file, v.Static, http.MethodGet)
	}
	v.route.Route("/", v.Static, http.MethodGet)

	return nil
}

//Down shutdown api service
func (v *Routes) Down() error {
	return nil
}

//Static controller
func (v *Routes) Static(w http.ResponseWriter, r *http.Request) {
	filename := r.RequestURI
	switch filename {
	case "", "/":
		filename = "/index.html"
		break
	}
	body := v.cache.Get(filename)

	contentType := http.DetectContentType(body)
	if contentType == "text/plain; charset=utf-8" {
		switch true {
		case strings.HasSuffix(filename, ".css"):
			contentType = "text/css; charset=utf-8"
			break
		case strings.HasSuffix(filename, ".js"):
			contentType = "application/javascript; charset=utf-8"
			break
		case strings.HasSuffix(filename, ".json"):
			contentType = "application/json; charset=utf-8"
			break
		}
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(body); err != nil {
		logger.Errorf("static response: %s", err.Error())
	}
}
