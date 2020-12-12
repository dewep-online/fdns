/*
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package rules

import "regexp"

type (
	RuleMap map[string]*Rule
	Rule    struct {
		reg *regexp.Regexp
		ip4 string
		ip6 string
	}
)

func (o *Rule) Compile(name string) (string, string, bool) {
	matches := o.reg.FindStringSubmatchIndex(name)
	if matches == nil {
		return "", "", false
	}
	ipv4 := o.reg.ExpandString([]byte{}, o.ip4, name, matches)
	ipv6 := o.reg.ExpandString([]byte{}, o.ip6, name, matches)
	return string(ipv4), string(ipv6), true
}

func (o *Rule) Match(name string) (string, string, bool) {
	if !o.reg.MatchString(name) {
		return "", "", false
	}
	return o.ip4, o.ip6, true
}
