package webserver

import "github.com/deweppro/go-http/servers"

type (
	BaseConfig struct {
		Middleware Middleware `yaml:"middleware"`
	}
	Middleware struct {
		Throttling int64 `yaml:"throttling"`
	}
)

type WebConfig struct {
	Http  servers.Config `yaml:"http"`
	Debug servers.Config `yaml:"debug"`
}
