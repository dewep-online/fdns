package pkg

import (
	"github.com/dewep-online/fdns/pkg/blacklist"
	"github.com/dewep-online/fdns/pkg/cache"
	"github.com/dewep-online/fdns/pkg/database"
	"github.com/dewep-online/fdns/pkg/dnscli"
	"github.com/dewep-online/fdns/pkg/rules"
	"github.com/deweppro/go-app/application"
	"github.com/deweppro/go-orm/schema/sqlite"
)

var (
	//Module di injector
	Module = application.Modules{
		blacklist.New,
		cache.New,
		dnscli.New,
		rules.New,
		database.New,
	}
	//Config di injector
	Config = application.Modules{
		&rules.Config{},
		&sqlite.Config{},
	}
)
