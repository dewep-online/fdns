/*
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package dnssrv

import "github.com/deweppro/go-app"

var (
	Module = app.Modules{
		New,
	}
	Config = app.Modules{
		&ConfigTCP{},
	}
)
