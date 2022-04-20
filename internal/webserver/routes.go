package webserver

import (
	"net/http"

	"github.com/deweppro/go-http/pkg/routes"
	"github.com/deweppro/go-logger"
	"github.com/deweppro/go-static"
)

//go:generate static ./../../web/dist/fdns ui

var ui static.Reader

//Routes model
type Routes struct {
	route *routes.Router
	conf  *BaseConfig
}

//NewRoutes init router
func New(c *BaseConfig, r *routes.Router) *Routes {
	return &Routes{
		route: r,
		conf:  c,
	}
}

//Up startup api service
func (v *Routes) Up() error {
	v.route.Global(routes.RecoveryMiddleware(logger.Default()))
	v.route.Global(routes.ThrottlingMiddleware(v.conf.Middleware.Throttling))

	for _, file := range ui.List() {
		logger.WithFields(logger.Fields{"url": file}).Infof("add static route")
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
	default:
	}

	if err := ui.ResponseWrite(w, filename); err != nil {
		logger.WithFields(logger.Fields{
			"err": err.Error(),
			"url": filename,
		}).Infof("static response")
	}
}
