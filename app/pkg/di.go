/*
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package pkg

import (
	"fdns/app/pkg/blacklist"
	"fdns/app/pkg/cache"
	"fdns/app/pkg/dnscli"
	"fdns/app/pkg/rules"

	"github.com/deweppro/go-app"
)

var (
	Module = app.Modules{
		blacklist.New,
		cache.New,
		dnscli.New,
		rules.New,
	}
	Config = app.Modules{
		&blacklist.Config{},
		&dnscli.Config{},
		&rules.Config{},
	}
)
