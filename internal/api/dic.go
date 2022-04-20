package api

import (
	"github.com/deweppro/go-app/application"
)

var (
	//Module di injector
	Module = application.Modules{
		NewAPI,
	}
	//Config di injector
	Config = application.Modules{}
)
