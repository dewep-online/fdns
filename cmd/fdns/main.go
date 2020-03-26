/*
 * Copyright (c) 2020.  Mikhail Knyazhev <markus621@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/gpl-3.0.html>.
 */

package main

import (
	"flag"
	"runtime"

	dns_client "fdns/internal/dns-client"

	"github.com/deweppro/core/pkg/app"

	app2 "fdns/internal/app"
	"fdns/internal/server"
)

var cfile = flag.String("config", "./configs/config.yaml", "path to config file")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	app.New(
		*cfile,
		app.NewInterfaces().Add(
			&app2.ConfigApp{},
			&server.ConfigServer{},
		),
		app.NewInterfaces().Add(
			app2.New,
			server.New,
			dns_client.New,
			server.Services{},
		),
	).Run()
}
