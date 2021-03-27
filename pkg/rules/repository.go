package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dewep-games/fdns/pkg/blacklist"
	"github.com/dewep-games/fdns/pkg/cache"
	"github.com/dewep-games/fdns/pkg/dnscli"
	"github.com/dewep-games/fdns/pkg/utils"
	"github.com/deweppro/go-logger"
	"github.com/miekg/dns"
)

type Repository struct {
	conf      *Config
	cache     *cache.Repository
	cli       *dnscli.Client
	blacklist *blacklist.Repository
	resolve   RuleMap
	regexp    RuleMap
}

func New(c *Config, r *cache.Repository, d *dnscli.Client, b *blacklist.Repository) *Repository {
	return &Repository{
		conf:      c,
		cache:     r,
		cli:       d,
		blacklist: b,
		regexp:    make(RuleMap),
		resolve:   make(RuleMap),
	}
}

func (o *Repository) Up() error {
	for domain, ips := range o.conf.HostRules {
		ip4, ip6 := utils.ParseIPs(ips)
		o.cache.Set(
			fmt.Sprintf("%s.", domain),
			cache.Create(ip4, ip6, false),
		)
	}

	for domain, ips := range o.conf.DNSRules {
		ip4, ip6 := utils.ParseIPs(ips)

		domain = regexp.QuoteMeta(domain)
		domain = strings.ReplaceAll(domain, "\\?", "?")
		domain = strings.ReplaceAll(domain, "\\*", ".*")
		domain = fmt.Sprintf("^.*%s\\.$", strings.Trim(domain, "^$"))

		o.resolve[domain] = &Rule{
			reg: regexp.MustCompile(domain),
			ip4: ip4,
			ip6: ip6,
		}
	}

	for domain, ips := range o.conf.RegExpRules {

		domain = fmt.Sprintf("^%s\\.$", strings.Trim(domain, "^$"))
		o.regexp[domain] = &Rule{
			reg:    regexp.MustCompile(domain),
			format: ips,
		}
	}

	for domain, ips := range o.conf.QueryRules {
		ip4, ip6 := utils.ParseIPs(ips)

		domain = regexp.QuoteMeta(domain)
		domain = strings.ReplaceAll(domain, "\\?", ".")
		domain = strings.ReplaceAll(domain, "\\*", ".*")
		domain = fmt.Sprintf("^%s\\.$", strings.Trim(domain, "^$"))

		o.regexp[domain] = &Rule{
			reg: regexp.MustCompile(domain),
			ip4: ip4,
			ip6: ip6,
		}
	}

	return nil
}

func (o *Repository) Down() error {
	return nil
}

func (o *Repository) DNS(name string, m *dns.Msg) ([]dns.RR, error) {
	var (
		ips []string
		rm  *dns.Msg
		err error
	)
	for _, item := range o.resolve {
		ip4, ip6, ok := item.Match(name)
		if !ok {
			continue
		}
		if len(ip6) > 0 {
			ips = append(ips, ip6...)
		}
		if len(ip4) > 0 {
			ips = append(ips, ip4...)
		}
	}
	if len(ips) > 0 {
		rm, err = o.cli.Exchange(m, ips)
	} else {
		rm, err = o.cli.ExchangeRandomDNS(m)
	}
	if err != nil {
		return nil, err
	}
	return o.nslookup(name, rm), nil
}

func (o *Repository) GetA(name string) ([]dns.RR, error) {
	c := o.cache.Get(name)
	if c == nil {
		if c = o.compileRegexp(name); c == nil {
			return nil, utils.ErrCacheNotFound
		}
	}
	if !c.HasIP4() {
		return nil, utils.ErrEmptyIP
	}
	return utils.CreateA(name, c.GetIP4()), nil
}

func (o *Repository) GetAAAA(name string) ([]dns.RR, error) {
	c := o.cache.Get(name)
	if c == nil {
		if c = o.compileRegexp(name); c == nil {
			return nil, utils.ErrCacheNotFound
		}
	}
	if !c.HasIP6() {
		return nil, utils.ErrEmptyIP
	}
	return utils.CreateAAAA(name, c.GetIP6()), nil
}

func (o *Repository) compileRegexp(name string) (rip *cache.Item) {
	for reg, item := range o.regexp {
		if ip4, ip6, ok := item.Compile(name); ok {
			rip = cache.Create(ip4, ip6, false)
			o.cache.Set(name, rip)

			logger.Infof("rules match: DOMAIN: %s REGEXP: %s IPv4: %s IPv6: %s", name, reg, ip4, ip6)
			break
		}
	}
	return
}

func (o *Repository) nslookup(name string, rm *dns.Msg) (answer []dns.RR) {
	for _, answ := range rm.Answer {
		switch answ.(type) {
		case *dns.A:
			if o.blacklist.Has(answ.(*dns.A).A) {
				continue
			}
			o.cache.Update(name, []string{answ.(*dns.A).A.String()}, nil)
		case *dns.AAAA:
			if o.blacklist.Has(answ.(*dns.AAAA).AAAA) {
				continue
			}
			o.cache.Update(name, nil, []string{answ.(*dns.AAAA).AAAA.String()})
		}

		answer = append(answer, answ)
	}
	return
}
