package pkg

import (
	"github.com/dewep-games/fdns/pkg/blacklist"
	"github.com/dewep-games/fdns/pkg/cache"
	"github.com/dewep-games/fdns/pkg/dnscli"
	"github.com/dewep-games/fdns/pkg/rules"
	"github.com/deweppro/go-app"
)

var (
	//Module di injector
	Module = app.Modules{
		blacklist.New,
		cache.New,
		dnscli.New,
		rules.New,
	}
	//Config di injector
	Config = app.Modules{
		&blacklist.Config{},
		&dnscli.Config{},
		&rules.Config{},
	}
)
