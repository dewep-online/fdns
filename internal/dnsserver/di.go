package dnsserver

import "github.com/deweppro/go-app"

var (
	Module = app.Modules{
		New,
	}
	Config = app.Modules{
		&ConfigTCP{},
	}
)
