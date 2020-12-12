/*
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"flag"
	"runtime"

	"fdns/app"
)

var cfile = flag.String("config", "./config.yaml", "path to config file")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	app.Run(*cfile)
}
