/*
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package rules

import (
	"fmt"
	"regexp"
	"strings"

	"fdns/app/pkg/blacklist"
	"fdns/app/pkg/cache"
	"fdns/app/pkg/dnscli"
	"fdns/app/pkg/utils"

	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

const (
	typeNone   = "none"
	typeDNS    = "dns"
	typeRegexp = "regexp"
	typeQuery  = "query"
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
	var sub string
	for _, rule := range o.conf.Rules {
		switch rule.Type {
		case typeNone:
			o.cache.Set(
				fmt.Sprintf("%s.", rule.Rule),
				cache.Create([]string{rule.IP4}, []string{rule.IP6},
					false))
			continue

		case typeDNS:
			sub = regexp.QuoteMeta(rule.Rule)
			sub = strings.ReplaceAll(sub, "\\?", "?")
			sub = strings.ReplaceAll(sub, "\\*", ".*")

			if ip, er := utils.ValidateIP(rule.IP4); er == nil {
				rule.IP4 = ip
			} else {
				rule.IP4 = ""
			}

			if ip, er := utils.ValidateIP(rule.IP6); er == nil {
				rule.IP6 = ip
			} else {
				rule.IP6 = ""
			}

			sub = fmt.Sprintf("^.*%s\\.$", strings.Trim(sub, "^$"))
			o.resolve[sub] = &Rule{
				reg: regexp.MustCompile(sub),
				ip4: rule.IP4,
				ip6: rule.IP6,
			}

			continue

		case typeQuery:
			sub = regexp.QuoteMeta(rule.Rule)
			sub = strings.ReplaceAll(sub, "\\?", ".")
			sub = strings.ReplaceAll(sub, "\\*", ".*")

		case typeRegexp:
			sub = rule.Rule

		default:
			continue
		}

		sub = fmt.Sprintf("^%s\\.$", strings.Trim(sub, "^$"))
		o.regexp[sub] = &Rule{
			reg: regexp.MustCompile(sub),
			ip4: rule.IP4,
			ip6: rule.IP6,
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
			ips = append(ips, ip6)
		}
		if len(ip4) > 0 {
			ips = append(ips, ip4)
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
	var (
		ip4, ip6 string
		ok       bool
	)
	for reg, item := range o.regexp {
		if ip4, ip6, ok = item.Compile(name); ok {
			rip = cache.Create([]string{ip4}, []string{ip6}, false)
			o.cache.Set(name, rip)

			logrus.WithFields(logrus.Fields{
				"domain": name, "regexp": reg,
				"ip4": ip4, "ip6": ip6,
			}).Info("rules match")
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
			o.cache.Update(name, []string{answ.(*dns.A).A.String()}, []string{})
		case *dns.AAAA:
			if o.blacklist.Has(answ.(*dns.AAAA).AAAA) {
				continue
			}
			o.cache.Update(name, []string{}, []string{answ.(*dns.AAAA).AAAA.String()})
		}

		answer = append(answer, answ)
	}
	return
}
