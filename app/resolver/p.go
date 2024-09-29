package resolver

import (
	"github.com/osspkg/fdns/app/cache"
	"github.com/osspkg/fdns/app/dnscli"
	"github.com/osspkg/fdns/app/rules"
	"go.osspkg.com/goppy/v2/plugins"
)

var Plugin = plugins.Plugin{
	Inject: func(cli *dnscli.Client, cr *cache.Records, rr *rules.Rules) *Resolver {
		return NewResolver(cli, cr, rr)
	},
}
