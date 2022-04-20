package rules

import (
	"regexp"
	"sync"
	"time"

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
	resolve   []Resolver
	mux       sync.RWMutex
}

func New(c *Config, r *cache.Repository, d *dnscli.Client, b *blacklist.Repository) *Repository {
	return &Repository{
		conf:      c,
		cache:     r,
		cli:       d,
		blacklist: b,
		resolve:   make([]Resolver, 0),
	}
}

func (v *Repository) Up() error {
	if err := HostRules(v.conf.HostRules, v); err != nil {
		return err
	}
	if err := DNSRules(v.conf.DNSRules, v); err != nil {
		return err
	}
	if err := RegexpRules(v.conf.RegExpRules, v); err != nil {
		return err
	}
	if err := QueryRules(v.conf.QueryRules, v); err != nil {
		return err
	}

	time.AfterFunc(time.Minute, func() {
		AdblockRules(v.conf.AdblockRules, v)
	})

	return nil
}

func (v *Repository) Down() error {
	return nil
}

func (v *Repository) SetHostResolve(domain string, ip4, ip6 []string, ttl int64) {
	v.cache.Set(domain, ip4, ip6, ttl)
}

func (v *Repository) SetRexResolve(format string, rx *regexp.Regexp, ip4, ip6 []string, tp uint) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.resolve = append(v.resolve, Resolver{
		reg:    rx,
		format: format,
		types:  tp,
		ip4:    ip4,
		ip6:    ip6,
	})
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
			if ttl > 0 {
				ttl = time.Now().Add(time.Second * time.Duration(ttl)).Unix()
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
	for _, vv := range resp.Answer {
		switch vv := vv.(type) {
		case *dns.A:
			if v.blacklist.Has(vv.A) {
				continue
			}
			ip = append(ip, vv.A.String())
		case *dns.AAAA:
			if v.blacklist.Has(vv.AAAA) {
				continue
			}
			ip = append(ip, vv.AAAA.String())
		}

		if ttl < vv.Header().Ttl {
			ttl = vv.Header().Ttl
		}
	}

	return time.Now().Add(time.Second * time.Duration(ttl)).Unix(), ip
}
