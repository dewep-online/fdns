/*
 *  Copyright (c) 2020-2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package rules

type Rules struct {
	adblock *AdBlock
	regexp  *RegexpRules
	static  *StaticRules
}

func NewRules(ab *AdBlock, rr *RegexpRules, sr *StaticRules) (*Rules, error) {
	return &Rules{
		adblock: ab,
		regexp:  rr,
		static:  sr,
	}, nil
}

func (v *Rules) IsBlocked(domain string) bool {
	return v.adblock.Contain(domain)
}

func (v *Rules) Static(qtype uint16, domain string) []string {
	if value := v.static.Convert(qtype, domain); len(value) > 0 {
		return value
	}
	return v.regexp.Convert(qtype, domain)
}
