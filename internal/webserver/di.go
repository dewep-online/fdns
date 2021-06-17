package webserver

import (
	"github.com/deweppro/go-app/application"
	"github.com/deweppro/go-http/web/server"
)

var (
	//Module di injector
	Module = application.Modules{
		server.New,
		NewRoutes,
	}
	//Config di injector
	Config = application.Modules{
		&server.Config{},
		&MiddlewareConfig{},
	}
)
