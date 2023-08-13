package main

import (
	"github.com/osspkg/fdns/app/cache"
	"github.com/osspkg/fdns/app/db"
	"github.com/osspkg/fdns/app/dns"
	"github.com/osspkg/fdns/app/dnscli"
	"github.com/osspkg/fdns/app/httpcli"
	"github.com/osspkg/fdns/app/resolver"
	"github.com/osspkg/fdns/app/rules"
	"github.com/osspkg/goppy"
	"github.com/osspkg/goppy/plugins/database"
	"github.com/osspkg/goppy/plugins/web"
)

func main() {
	app := goppy.New()
	app.Plugins(
		web.WithHTTP(),
		database.WithMySQL(),
	)
	app.Plugins(cache.Plugins...)
	app.Plugins(rules.Plugins...)
	app.Plugins(
		db.Plugin,
		httpcli.Plugin,
		resolver.Plugin,
		dns.Plugin,
		dnscli.Plugin,
	)
	app.Run()
}
