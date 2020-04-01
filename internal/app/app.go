/*
 * Copyright (c) 2020.  Mikhail Knyazhev <markus621@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/gpl-3.0.html>.
 */

package app

import (
	"net"

	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"

	"fdns/internal/utils"
)

const (
	cTTL = 900
)

type App struct {
	config *ConfigApp

	cache  *DNSCache
	regexp RegexpRules
	regdns RegexpRules

	blacklistIP    []net.IP
	blacklistIPNet []*net.IPNet
}

func New(c *ConfigApp) *App {
	app := &App{
		config:         c,
		cache:          NewDNSCache(c),
		regexp:         make(RegexpRules),
		regdns:         make(RegexpRules),
		blacklistIP:    make([]net.IP, 0),
		blacklistIPNet: make([]*net.IPNet, 0),
	}

	app.blacklist()
	app.rules()

	return app
}

func (a *App) GetDNS(name string) (result []string, err error) {
	for _, item := range a.regdns {
		ip4, ip6, ok := item.Match(name)
		if !ok {
			continue
		}
		if len(ip6) > 0 {
			result = append(result, ip6)
		}
		if len(ip4) > 0 {
			result = append(result, ip4)
		}
	}
	if len(result) == 0 {
		err = utils.ErrorEmptyIP
	}
	return
}

func (a *App) GetA(name string) (dns.RR, error) {
	c := a.cache.Get(name)
	if c == nil {
		if c = a.compileRegexp(name); c == nil {
			return nil, utils.ErrorCacheNotFound
		}
	}

	if c.IsEmptyIP4() {
		if c.IsUpdatable() {
			return nil, utils.ErrorCacheNotFound
		}
		return nil, utils.ErrorEmptyIP
	}

	return a.makeA(name, c.GetIP4()), nil
}

func (a *App) GetAAAA(name string) (dns.RR, error) {
	c := a.cache.Get(name)
	if c == nil {
		if c = a.compileRegexp(name); c == nil {
			return nil, utils.ErrorCacheNotFound
		}
	}

	if c.IsEmptyIP6() {
		if c.IsUpdatable() {
			return nil, utils.ErrorCacheNotFound
		}
		return nil, utils.ErrorEmptyIP
	}

	return a.makeAAAA(name, c.GetIP6()), nil
}

func (a *App) makeA(name, ip string) dns.RR {
	return &dns.A{
		Hdr: dns.RR_Header{
			Name:   name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    uint32(cTTL),
		},
		A: net.ParseIP(ip),
	}
}

func (a *App) makeAAAA(name, ip string) dns.RR {
	return &dns.AAAA{
		Hdr: dns.RR_Header{
			Name:   name,
			Rrtype: dns.TypeAAAA,
			Class:  dns.ClassINET,
			Ttl:    uint32(cTTL),
		},
		AAAA: net.ParseIP(ip),
	}
}

func (a *App) compileRegexp(name string) (rip *RuleIP) {
	var (
		ip4, ip6 string
		ok       bool
	)

	for reg, item := range a.regexp {
		if ip4, ip6, ok = item.Compile(name); ok {

			rip = NewRuleIP(ip4, ip6, false)
			a.cache.Set(name, rip)

			logrus.WithFields(logrus.Fields{
				"domain": name, "regexp": reg,
				"ip4": ip4, "ip6": ip6,
			}).Info("rules match")

			break
		}
	}
	return
}

func (a *App) Update(name, ip4, ip6 string) {
	a.cache.Update(name, ip4, ip6)
}
