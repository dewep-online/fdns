package pkg

import (
	"github.com/dewep-online/fdns/pkg/blacklist"
	"github.com/dewep-online/fdns/pkg/cache"
	"github.com/dewep-online/fdns/pkg/dnscli"
	"github.com/dewep-online/fdns/pkg/rules"
	"github.com/deweppro/go-app/application"
)

var (
	//Module di injector
	Module = application.Modules{
		blacklist.New,
		cache.New,
		dnscli.New,
		rules.New,
	}
	//Config di injector
	Config = application.Modules{
		&blacklist.Config{},
		&dnscli.Config{},
		&rules.Config{},
	}
)
