/*
 *  Copyright (c) 2020-2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package resolver

import (
	"github.com/miekg/dns"
	"github.com/osspkg/fdns/app/cache"
	"github.com/osspkg/fdns/app/dnscli"
	"github.com/osspkg/fdns/app/rules"
	"go.osspkg.com/logx"
)

type (
	Resolver struct {
		cli   *dnscli.Client
		cache *cache.Records
		rules *rules.Rules
	}
)

func NewResolver(c *dnscli.Client, r *cache.Records, br *rules.Rules) *Resolver {
	return &Resolver{
		cli:   c,
		cache: r,
		rules: br,
	}
}

func (v *Resolver) Resolve(question dns.Question) []dns.RR {
	if c, ok := v.cache.Get(question.Qtype, question.Name); ok {
		return CreateRR(question.Qtype, question.Name, c.TTL, c.Value...)
	}

	if v.rules.IsBlocked(question.Name) {
		return nil
	}

	if value := v.rules.Static(question.Qtype, question.Name); len(value) > 0 {
		return CreateRR(question.Qtype, question.Name, DefaultTTL, value...)
	}

	response, err := v.cli.Exchange(question)
	if err != nil {
		logx.Error("DNS Exchange", "err", err, "domain", question.Name)
		return nil
	}

	value, ttl := ParseRR(response)
	if len(value) > 0 {
		v.cache.Set(question.Qtype, question.Name, ttl, value...)
	}

	return response
}
