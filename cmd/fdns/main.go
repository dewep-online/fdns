package main

import (
	"flag"

	"github.com/dewep-online/fdns/internal/api"
	"github.com/dewep-online/fdns/internal/dnsserver"
	"github.com/dewep-online/fdns/internal/webserver"
	"github.com/dewep-online/fdns/pkg"
	"github.com/deweppro/go-app/application"
	"github.com/deweppro/go-http/web/debug"
	"github.com/deweppro/go-logger"
)

var configFile = flag.String("config", "./config.yaml", "path to config file")

func main() {
	flag.Parse()

	application.New().
		Logger(logger.Default()).
		ConfigFile(
			*configFile,
			&debug.Config{},
			pkg.Config,
			webserver.Config,
			dnsserver.Config,
			api.Config,
		).
		Modules(
			debug.New,
			pkg.Module,
			webserver.Module,
			dnsserver.Module,
			api.Module,
		).
		Run()
}
