package dnsserver

import "github.com/deweppro/go-app/application"

var (
	Module = application.Modules{
		New,
	}
	Config = application.Modules{
		&ConfigTCP{},
	}
)
