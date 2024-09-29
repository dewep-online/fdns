/*
 *  Copyright (c) 2020-2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package dnscli

import (
	"github.com/osspkg/fdns/app/db"
	"go.osspkg.com/goppy/v2/plugins"
	"go.osspkg.com/goppy/v2/xdns"
)

var Plugin = plugins.Plugin{
	Inject: func(dbc db.Connect, xs *xdns.Server, xc *xdns.Client) *Client {
		cli := NewClient(dbc, xc)
		xs.HandleFunc(xc)
		return cli
	},
}
