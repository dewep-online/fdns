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
	"fmt"
	"regexp"
	"strings"
)

const (
	cTypeNone   = "none"
	cTypeRegexp = "regexp"
	cTypeQuery  = "query"
)

type RegexpRules map[string]*RegexpRule
type RegexpRule struct {
	reg *regexp.Regexp
	ip4 string
	ip6 string
}

func (rx *RegexpRule) Compile(name string) (string, string, bool) {
	matches := rx.reg.FindStringSubmatchIndex(name)
	if matches == nil {
		return "", "", false
	}

	ipv4 := rx.reg.ExpandString([]byte{}, rx.ip4, name, matches)
	ipv6 := rx.reg.ExpandString([]byte{}, rx.ip6, name, matches)

	return string(ipv4), string(ipv6), true
}

func (a *App) rules() {
	var sub string
	for _, rule := range a.config.Rules {
		switch rule.Type {
		case cTypeNone:
			a.cache.Set(fmt.Sprintf("%s.", rule.Rule), NewRuleIP(rule.IP4, rule.IP6, false))
			continue

		case cTypeQuery:
			sub = regexp.QuoteMeta(rule.Rule)
			sub = strings.ReplaceAll(sub, "\\?", "?")
			sub = strings.ReplaceAll(sub, "\\*", ".*")

		case cTypeRegexp:
			sub = rule.Rule

		default:
			continue
		}

		sub = fmt.Sprintf("^%s\\.$", strings.Trim(sub, "^$"))
		a.regexp[sub] = &RegexpRule{
			reg: regexp.MustCompile(sub),
			ip4: rule.IP4,
			ip6: rule.IP6,
		}
	}
}
