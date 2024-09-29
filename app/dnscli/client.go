/*
 *  Copyright (c) 2020-2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package dnscli

import (
	"context"
	"encoding/json"
	"time"

	"github.com/miekg/dns"
	"github.com/osspkg/fdns/app/db"
	"go.osspkg.com/goppy/v2/orm"
	"go.osspkg.com/goppy/v2/xdns"
	"go.osspkg.com/logx"
	"go.osspkg.com/network/address"
	"go.osspkg.com/random"
	"go.osspkg.com/routine"
	"go.osspkg.com/syncing"
	"go.osspkg.com/validate"
	"go.osspkg.com/xc"
)

type (
	Rules struct {
		data map[string][]string
	}
)

func NewRules() *Rules {
	return &Rules{
		data: make(map[string][]string, 100),
	}
}

func (v *Rules) Set(zone string, ips []string) {
	_, ok := v.data[zone]
	if !ok {
		v.data[zone] = make([]string, 0, 2)
	}

	v.data[zone] = append(v.data[zone], ips...)
}

func (v *Rules) Resolve(zone string) (result []string) {
	for i := 2; i >= 0; i-- {
		vv := validate.GetDomainLevel(zone, i)
		if ips, ok := v.data[vv]; ok {
			result = append(result, ips...)
			return random.Shuffle(result)
		}
	}
	return
}

// --------------------------------------------------------------------------------------------------------------------

type (
	Client struct {
		ns  *Rules
		cli *xdns.Client
		db  db.Connect
		mux syncing.Lock
	}
)

func NewClient(dbc db.Connect, cli *xdns.Client) *Client {
	c := &Client{
		cli: cli,
		ns:  NewRules(),
		db:  dbc,
		mux: syncing.NewLock(),
	}

	c.cli.SetZoneResolver(c.ns)

	return c
}

func (v *Client) Up(ctx xc.Context) error {
	routine.Interval(
		ctx.Context(),
		15*time.Minute,
		func(ctx context.Context) {
			if err := v.ForceUpdate(ctx); err != nil {
				logx.Error("DNS Client update dns list", "err", err)
			}
		},
	)
	return nil
}

func (v *Client) Down() error {
	return nil
}

func (v *Client) ForceUpdate(ctx context.Context) error {
	ns := NewRules()
	err := v.db.Main().Query(ctx, "load_ns_zone", func(q orm.Querier) {
		q.SQL("SELECT `zone`,`data` FROM `dns`;")
		q.Bind(func(bind orm.Scanner) error {
			var (
				zone string
				b    []byte
			)
			if err := bind.Scan(&zone, &b); err != nil {
				return err
			}
			var data []string
			if err := json.Unmarshal(b, &data); err != nil {
				return err
			}
			ns.Set(zone, address.Normalize("53", data...))
			return nil
		})
	})
	if err != nil {
		return err
	}
	v.mux.Lock(func() {
		v.ns = ns
	})
	return nil
}

func (v *Client) Exchange(question dns.Question) ([]dns.RR, error) {
	return v.cli.Exchange(question)
}
