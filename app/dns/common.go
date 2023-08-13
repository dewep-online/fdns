package dns

import (
	"github.com/osspkg/fdns/app/resolver"
	"github.com/osspkg/goppy/plugins"
)

var Plugin = plugins.Plugin{
	Inject: func(c *Config, r resolver.Resolver) *Server {
		return NewServer(c.DNS, r)
	},
	Config: &Config{},
}
