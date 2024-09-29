/*
 *  Copyright (c) 2020-2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package main

import (
	"github.com/osspkg/fdns/app/api"
	"github.com/osspkg/fdns/app/cache"
	"github.com/osspkg/fdns/app/db"
	"github.com/osspkg/fdns/app/dnscli"
	"github.com/osspkg/fdns/app/resolver"
	"github.com/osspkg/fdns/app/rules"
	"go.osspkg.com/goppy/v2"
	"go.osspkg.com/goppy/v2/orm"
	"go.osspkg.com/goppy/v2/web"
	"go.osspkg.com/goppy/v2/xdns"
)

func main() {
	app := goppy.New("fdns", "v0.0.0-dev", "filter dns")
	app.Plugins(
		web.WithServer(),
		web.WithClient(),
		orm.WithMysql(),
		orm.WithORM(),
		xdns.WithServer(),
		xdns.WithClient(),
	)
	app.Plugins(api.Plugins...)
	app.Plugins(cache.Plugins...)
	app.Plugins(rules.Plugins...)
	app.Plugins(
		db.Plugin,
		resolver.Plugin,
		dnscli.Plugin,
	)
	app.Run()
}
