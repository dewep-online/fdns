/*
 *  Copyright (c) 2020-2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package api

import (
	"github.com/osspkg/fdns/app/db"
	"go.osspkg.com/goppy/v2/web"
)

type Api struct {
	router web.Router
	db     db.Connect
}

func NewApi(r web.RouterPool, dbc db.Connect) *Api {
	return &Api{
		router: r.Main(),
		db:     dbc,
	}
}

func (v *Api) Up() error {
	api := v.router.Collection("/api")
	api.Get("/blacklist/adblock/list", v.AdblockList)
	return nil
}

func (v *Api) Down() error {
	return nil
}
