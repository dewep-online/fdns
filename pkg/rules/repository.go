package rules

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/deweppro/go-errors"

	"github.com/deweppro/go-logger"

	"github.com/dewep-online/fdns/pkg/database"

	"github.com/deweppro/go-app/application/ctx"

	"github.com/dewep-online/fdns/pkg/utils"

	"github.com/dewep-online/fdns/pkg/blacklist"
	"github.com/dewep-online/fdns/pkg/cache"
	"github.com/dewep-online/fdns/pkg/dnscli"
	"github.com/miekg/dns"
)

type Repository struct {
	conf      *Config
	cache     *cache.Repository
	cli       *dnscli.Client
	blacklist *blacklist.Repository
	db        *database.Database
	resolve   map[string]Resolver
	mux       sync.RWMutex
}

func New(c *Config, r *cache.Repository, d *dnscli.Client, b *blacklist.Repository, db *database.Database) *Repository {
	return &Repository{
		conf:      c,
		cache:     r,
		cli:       d,
		blacklist: b,
		db:        db,
		resolve:   make(map[string]Resolver, 0),
	}
}

func (v *Repository) Up(ctx ctx.Context) error {
	var timestamp int64

	utils.Interval(ctx.Context(), time.Hour*24, func(ctx context.Context) {
		uris := LoadAdblockRules(v.conf.AdblockRules)
		AdblockRules(uris, func(uri string, domains []string) {
			tag, err := v.db.SetBlacklistURI(ctx, uri)
			if err != nil {
				logger.Warnf("adblock-uri [%s]: %s", uri, utils.StringError(err))
				return
			}
			if err = v.db.SetBlacklistDomain(ctx, tag, domains); err != nil {
				logger.Warnf("adblock-uri [%s]: %s", uri, utils.StringError(err))
				return
			}
		})
	})

	utils.Interval(ctx.Context(), time.Minute*5, func(ctx context.Context) {
		err := errors.Wrap(
			v.db.GetRulesMap(ctx, database.DNS, timestamp, func(m map[string]string) error {
				return DNSRules(m, v)
			}),
			v.db.GetRulesMap(ctx, database.Host, timestamp, func(m map[string]string) error {
				return HostRules(m, v)
			}),
			v.db.GetRulesMap(ctx, database.Regex, timestamp, func(m map[string]string) error {
				return RegexpRules(m, v)
			}),
			v.db.GetRulesMap(ctx, database.Query, timestamp, func(m map[string]string) error {
				return QueryRules(m, v)
			}),
			func() error {
				rules, err := v.db.GetBlacklistDomain(ctx, timestamp)
				if err != nil {
					return err
				}
				logger.Infof("update rules [adblock]: %d", len(rules))
				return HostRules(rules.ToMap(database.ActiveTrue), v)
			}(),
		)
		if err != nil {
			logger.Warnf("update rules: %s", utils.StringError(err))
		}

		timestamp = time.Now().Unix()
	})

	return nil
}

func (v *Repository) Down(ctx ctx.Context) error {
	return nil
}

func (v *Repository) SetHostResolve(domain string, ip4, ip6 []string, ttl int64) {
	v.cache.Set(domain, ip4, ip6, ttl)
}

func (v *Repository) ReplaceRexResolve(t database.Types, o, n, ips string) {
	m := map[string]string{n: ips}
	switch t {
	case database.DNS:
		_ = DNSRules(m, v)
	case database.Regex:
		_ = RegexpRules(m, v)
	case database.Query:
		_ = QueryRules(m, v)
	}

	v.mux.Lock()
	r, ok := v.resolve[o]
	if o != n {
		delete(v.resolve, o)
	}
	v.mux.Unlock()

	if ok {
		v.cache.DelByCallback(func(name string) bool {
			_, _, ok = r.Match(name)
			fmt.Println(r.rule, name, ok)
			return ok
		})
	}
}

func (v *Repository) DeleteRexResolve(name string) {
	v.mux.Lock()
	r, ok := v.resolve[name]
	delete(v.resolve, name)
	v.mux.Unlock()

	if ok {
		v.cache.DelByCallback(func(name string) bool {
			_, _, ok = r.Match(name)
			return ok
		})
	}
}

func (v *Repository) SetRexResolve(rule, format string, rx *regexp.Regexp, ip4, ip6 []string, tp uint) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.resolve[rule] = Resolver{
		rule:   rule,
		reg:    rx,
		format: format,
		types:  tp,
		ip4:    ip4,
		ip6:    ip6,
	}
}

func (v *Repository) Resolve(q dns.Question) []dns.RR {
	var (
		tp        uint
		ttl, ttl6 int64
		ip4, ip6  []string
	)

	tp, ip4, ip6 = v.rxlookup(q.Name)
	if tp == TypeNone || tp == TypeDNS {
		ips := append(ip4, ip6...)
		ttl, ip4 = v.nslookup(new(dns.Msg).SetQuestion(dns.Fqdn(q.Name), dns.TypeA), ips)
		ttl6, ip6 = v.nslookup(new(dns.Msg).SetQuestion(dns.Fqdn(q.Name), dns.TypeAAAA), ips)

		if len(ip4) > 0 || len(ip6) > 0 {
			if ttl6 > ttl {
				ttl = ttl6
			}
			v.cache.Set(q.Name, ip4, ip6, ttl)
		}
	}

	switch q.Qtype {
	case dns.TypeA:
		return utils.CreateA(q.Name, ip4)
	case dns.TypeAAAA:
		return utils.CreateAAAA(q.Name, ip6)
	}

	return nil
}

func (v *Repository) rxlookup(domain string) (uint, []string, []string) {
	if vv := v.cache.Get(domain); vv != nil {
		return TypeHost, vv.GetIP4(), vv.GetIP6()
	}
	v.mux.RLock()
	defer v.mux.RUnlock()

	for _, item := range v.resolve {
		ip4, ip6, ok := item.Compile(domain)
		if !ok {
			ip4, ip6, ok = item.Match(domain)
		}
		if !ok {
			continue
		}
		if item.types == TypeRegexp {
			v.cache.Set(domain, ip4, ip6, 0)
		}
		return item.types, ip4, ip6
	}
	return TypeNone, nil, nil
}

func (v *Repository) nslookup(msg *dns.Msg, ips []string) (int64, []string) {
	var (
		resp *dns.Msg
		err  error
	)
	if len(ips) > 0 {
		resp, err = v.cli.Exchange(msg, ips)
	} else {
		resp, err = v.cli.ExchangeRandomDNS(msg)
	}
	if err != nil {
		return 0, nil
	}

	var (
		ip  []string
		ttl uint32
	)
	if resp == nil {
		return time.Now().Add(time.Second * time.Duration(ttl)).Unix(), ip
	}
	for _, vv := range resp.Answer {
		switch vv.(type) {
		case *dns.A:
			if v.blacklist.Has(vv.(*dns.A).A) {
				continue
			}
			ip = append(ip, vv.(*dns.A).A.String())
		case *dns.AAAA:
			if v.blacklist.Has(vv.(*dns.AAAA).AAAA) {
				continue
			}
			ip = append(ip, vv.(*dns.AAAA).AAAA.String())
		}

		if ttl < vv.Header().Ttl {
			ttl = vv.Header().Ttl
		}
	}

	return time.Now().Add(time.Second * time.Duration(ttl)).Unix(), ip
}
