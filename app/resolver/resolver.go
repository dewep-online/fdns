package resolver

import (
	"github.com/miekg/dns"
	"github.com/osspkg/fdns/app/cache"
	"github.com/osspkg/fdns/app/dnscli"
	"github.com/osspkg/fdns/app/rules"
	"github.com/osspkg/goppy/plugins"
)

var Plugin = plugins.Plugin{
	Inject: func(c dnscli.Client, r *cache.Records, br *rules.Rules) Resolver {
		return NewResolver(c, r, br)
	},
}

type (
	Resolver interface {
		Resolve(q dns.Question) []dns.RR
	}
	object struct {
		cli   dnscli.Client
		cache *cache.Records
		rules *rules.Rules
	}
)

func NewResolver(c dnscli.Client, r *cache.Records, br *rules.Rules) Resolver {
	return &object{
		cli:   c,
		cache: r,
		rules: br,
	}
}

func (v *object) Resolve(question dns.Question) []dns.RR {
	if response := v.cache.Get(question.Qtype, question.Name); response != nil {
		return CreateRR(question.Qtype, question.Name, response.Ttl(), response.Values()...)
	}

	if v.rules.IsBlocked(question.Name) {
		return nil
	}

	response, ok := v.cli.Exchange(question)
	if !ok {
		return nil
	}

	value, ttl := ParseRR(response)
	v.cache.Set(question.Qtype, question.Name, ttl, value...)

	return response
}
