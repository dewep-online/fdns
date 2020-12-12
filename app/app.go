/*
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

import (
	"fdns/app/dnssrv"
	"fdns/app/httpsrv"
	"fdns/app/pkg"

	"github.com/sirupsen/logrus"

	application "github.com/deweppro/go-app"
)

func Run(conf string) {
	application.New().
		Logger(logrus.StandardLogger()).
		ConfigFile(
			conf,
			dnssrv.Config,
			httpsrv.Config,
			pkg.Config,
		).
		Modules(
			dnssrv.Module,
			httpsrv.Module,
			pkg.Module,
		).
		Run()
}
