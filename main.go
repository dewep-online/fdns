package main

import (
	"fdns/internal"
	"flag"
)

var cfile = flag.String("config", "./config.yaml", "path to config file")

func main() {
	flag.Parse()
	internal.Run(*cfile)
}
