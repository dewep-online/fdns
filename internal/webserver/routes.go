package webserver

import (
	"net/http"

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

	if err := v.cache.Write(filename, w); err != nil {
		logger.Errorf("static response: %s", err.Error())
	}
}
