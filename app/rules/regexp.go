/*
 *  Copyright (c) 2020-2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package rules

import (
	"context"
	"encoding/json"
	"regexp"
	"sync"
	"time"

	"github.com/osspkg/fdns/app/db"
	"go.osspkg.com/goppy/v2/orm"
	"go.osspkg.com/logx"
	"go.osspkg.com/routine"
	"go.osspkg.com/xc"
)

type RegexpRule struct {
	rule  *regexp.Regexp
	qtype uint16
	data  []string
}

func NewRegexpRule(rule string, qtype uint16, data []string) (*RegexpRule, error) {
	rx, err := regexp.Compile(rule)
	if err != nil {
		return nil, err
	}
	return &RegexpRule{
		rule:  rx,
		qtype: qtype,
		data:  data,
	}, nil
}

func (v *RegexpRule) Compile(name string) []string {
	if v.rule == nil {
		return nil
	}
	result := make([]string, 0, len(v.data))
	matches := v.rule.FindStringSubmatchIndex(name)
	if matches == nil {
		return append(result, v.data...)
	}
	for _, s := range v.data {
		value := v.rule.ExpandString([]byte{}, s, name, matches)
		result = append(result, string(value))
	}
	return result
}

func (v *RegexpRule) Match(name string) bool {
	if v.rule == nil {
		return false
	}
	return v.rule.MatchString(name)
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type RegexpRules struct {
	db   db.Connect
	data map[uint16][]*RegexpRule
	mux  sync.RWMutex
}

func NewRegexpRules(dbc db.Connect) *RegexpRules {
	return &RegexpRules{
		db:   dbc,
		data: make(map[uint16][]*RegexpRule, 100),
	}
}

func (v *RegexpRules) Up(ctx xc.Context) error {
	routine.Interval(ctx.Context(), time.Hour, func(ctx context.Context) {
		if err := v.ForceUpdate(ctx); err != nil {
			logx.Error("RegexpRules update", "err", err)
		}
	})
	return nil
}

func (v *RegexpRules) Down() error {
	return nil
}

func (v *RegexpRules) ForceUpdate(ctx context.Context) error {
	result := make(map[uint16][]*RegexpRule, 100)
	err := v.db.Main().Query(ctx, "load_regexp_rules", func(q orm.Querier) {
		q.SQL("SELECT `rule`, `qtype`,`data` FROM `static_regexp_rules` WHERE `deleted_at` IS NULL;")
		q.Bind(func(bind orm.Scanner) error {
			var (
				rule  string
				qtype uint16
				b     []byte
				data  []string
			)
			if err := bind.Scan(&rule, &qtype, &b); err != nil {
				return err
			}
			if err := json.Unmarshal(b, &data); err != nil {
				return err
			}
			rr, err := NewRegexpRule(rule, qtype, data)
			if err != nil {
				return err
			}
			if _, ok := result[qtype]; !ok {
				result[qtype] = make([]*RegexpRule, 0, 2)
			}
			result[qtype] = append(result[qtype], rr)
			return nil
		})
	})
	if err != nil {
		return err
	}
	v.mux.Lock()
	v.data = result
	v.mux.Unlock()
	return nil
}

func (v *RegexpRules) Convert(qtype uint16, domain string) []string {
	v.mux.RLock()
	defer v.mux.RUnlock()

	if _, ok := v.data[qtype]; !ok {
		return nil
	}

	for _, datum := range v.data[qtype] {
		if datum.Match(domain) {
			return datum.Compile(domain)
		}
	}
	return nil
}
