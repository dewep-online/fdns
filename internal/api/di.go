package api

import (
	"github.com/deweppro/go-app"
)

var (
	//Module di injector
	Module = app.Modules{
		NewAPI,
	}
	//Config di injector
	Config = app.Modules{}
)
