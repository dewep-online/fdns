package webserver

import (
	"github.com/deweppro/go-app"
	"github.com/deweppro/go-http/web/server"
)

var (
	//Module di injector
	Module = app.Modules{
		server.New,
		NewRoutes,
	}
	//Config di injector
	Config = app.Modules{
		&server.Config{},
		&MiddlewareConfig{},
	}
)
