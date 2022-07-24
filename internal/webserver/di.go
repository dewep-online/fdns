package webserver

import (
	"github.com/deweppro/go-app/application"
	"github.com/deweppro/go-http/pkg/routes"
	"github.com/deweppro/go-http/servers/debug"
	"github.com/deweppro/go-http/servers/web"
	"github.com/deweppro/go-logger"
)

var (
	//Module di injector
	Module = application.Modules{
		routes.NewRouter,
		WebAdapter,
		New,
	}
	//Config di injector
	Config = application.Modules{
		&WebConfig{},
		&MiddlewareConfig{},
	}
)

func WebAdapter(c *WebConfig, r *routes.Router) (*web.Server, *debug.Debug) {
	ws := web.New(c.Http, r, logger.Default())
	ds := debug.New(c.Debug, logger.Default())
	return ws, ds
}
