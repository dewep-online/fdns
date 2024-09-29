/*
 *  Copyright (c) 2020-2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package rules

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/osspkg/fdns/app/db"
	"go.osspkg.com/algorithms/filters/bloom"
	"go.osspkg.com/encrypt/hash"
	"go.osspkg.com/goppy/v2/orm"
	"go.osspkg.com/goppy/v2/web"
	"go.osspkg.com/logx"
	"go.osspkg.com/routine"
	"go.osspkg.com/syncing"
	"go.osspkg.com/validate"
	"go.osspkg.com/xc"
)

const (
	AdBlockStatic  = "static"
	AdBlockDynamic = "dynamic"
)

var rex = regexp.MustCompile(`(?miU)^\|\|([a-z0-9-.]+)\^(\n|\r|\$)`)

type AdBlock struct {
	bloom *bloom.Bloom
	cli   *web.ClientHttp
	db    db.Connect
	mux   syncing.Lock
}

func NewAdBlock(dbc db.Connect, cli web.ClientHttpPool) (*AdBlock, error) {
	ab := &AdBlock{
		cli: cli.Create(),
		db:  dbc,
		mux: syncing.NewLock(),
	}
	var err error
	if ab.bloom, err = ab.createBloom(10_000_000); err != nil {
		return nil, err
	}
	return ab, nil
}

func (v *AdBlock) createBloom(size uint64) (bf *bloom.Bloom, err error) {
	for i := 0; i < 10; i++ {
		bf, err = bloom.New(size*2, 0.0001)
		if err != nil {
			continue
		}
		return bf, nil
	}
	return
}

func (v *AdBlock) Up(ctx xc.Context) error {
	if err := v.ForceUpdate(ctx.Context()); err != nil {
		logx.Error("AdBlock update", "err", err)
	}
	go routine.Interval(ctx.Context(), time.Hour*24, func(ctx context.Context) {
		if err := v.UpgradeRules(ctx); err != nil {
			logx.Error("AdBlock update", "err", err)
		}
		if err := v.ForceUpdate(ctx); err != nil {
			logx.Error("AdBlock update", "err", err)
		}
	})
	return nil
}

func (v *AdBlock) Down() error {
	return nil
}

func (v *AdBlock) UpgradeRules(ctx context.Context) error {
	list := make(map[uint64]string, 10)
	err := v.db.Main().Query(ctx, "get_adblock_list", func(q orm.Querier) {
		q.SQL(
			"SELECT `id`,`data` FROM `blacklist_adblock_list` WHERE `deleted_at` IS NULL AND `type` = ?;",
			AdBlockDynamic,
		)
		q.Bind(func(bind orm.Scanner) error {
			var (
				id   uint64
				data string
			)
			if err := bind.Scan(&id, &data); err != nil {
				return err
			}
			list[id] = data
			return nil
		})
	})
	if err != nil {
		return err
	}

	for id, uri := range list {
		count, err := func() (int, error) {
			var b []byte
			err0 := v.cli.Call(ctx, http.MethodGet, uri, nil, &b)
			if err0 != nil {
				return 0, err0
			}

			rexResult := rex.FindAll(b, -1)
			result := make([]string, 0, 100)

			for _, rr := range rexResult {
				rule := strings.Trim(string(rr[2:len(rr)-1]), "\n^") + "."

				result = append(result, rule)
				if len(result) == 100 {
					if err0 = v.save(ctx, id, result); err0 != nil {
						return 0, err0
					}
					result = result[:0]
				}
			}

			if len(result) > 0 {
				if err0 = v.save(ctx, id, result); err0 != nil {
					return 0, err0
				}
			}

			return len(rexResult), nil
		}()
		if err != nil {
			logx.Error("AdBlock upgrade", "uri", uri, "err", err)
		} else {
			logx.Info("AdBlock upgrade", "uri", uri, "count", count)
		}
	}

	return nil
}

func (v *AdBlock) save(ctx context.Context, id uint64, data []string) error {
	return v.db.Main().Tx(ctx, "save_adblock_rules", func(v orm.Tx) {
		v.Exec(func(e orm.Executor) {
			e.SQL("INSERT IGNORE INTO `blacklist_adblock_rules` (`list_id`, `data`, `hash`, `updated_at`) VALUES (?, ?, ?, now());")
			for _, datum := range data {
				e.Params(id, datum, hash.SHA1(datum))
			}
		})
	})
}

func (v *AdBlock) ForceUpdate(ctx context.Context) error {
	count := 0
	err := v.db.Main().Query(ctx, "count_adblock_rules", func(q orm.Querier) {
		q.SQL("SELECT COUNT(*) FROM `blacklist_adblock_rules` WHERE `deleted_at` IS NULL;")
		q.Bind(func(bind orm.Scanner) error {
			return bind.Scan(&count)
		})
	})
	if err != nil {
		return err
	}
	if count <= 0 {
		return nil
	}

	var bf *bloom.Bloom
	if bf, err = v.createBloom(uint64(count)); err != nil {
		return err
	}

	err = v.db.Main().Query(ctx, "load_adblock_rules", func(q orm.Querier) {
		q.SQL("SELECT `data` FROM `blacklist_adblock_rules` WHERE `deleted_at` IS NULL;")
		q.Bind(func(bind orm.Scanner) error {
			var data string
			if err0 := bind.Scan(&data); err0 != nil {
				return err0
			}
			bf.Add([]byte(data))
			return nil
		})
	})
	if err != nil {
		return err
	}

	v.mux.Lock(func() {
		v.bloom = bf
	})

	return nil
}

func (v *AdBlock) Contain(name string) bool {
	has := false
	levels := validate.CountDomainLevels(name)
	params := make([]interface{}, 0, levels)
	for i := validate.CountDomainLevels(name); i >= 1; i-- {
		subDomain := validate.GetDomainLevel(name, i)
		params = append(params, subDomain)
		v.mux.RLock(func() {
			has = has || v.bloom.Contain([]byte(subDomain))
		})
	}
	if !has || len(params) == 0 {
		return false
	}

	count := 0
	err := v.db.Main().Query(context.Background(), "get_one_adblock_rule", func(q orm.Querier) {
		q.SQL(
			fmt.Sprintf(
				"SELECT COUNT(*) FROM `blacklist_adblock_rules` WHERE `deleted_at` IS NULL AND `data` IN (%s);",
				strings.Trim(strings.Repeat("?,", len(params)), ","),
			),
			params...,
		)
		q.Bind(func(bind orm.Scanner) error {
			return bind.Scan(&count)
		})
	})
	if err != nil {
		logx.Error("Check domain in adblock", "err", err, "domain", name)
	}

	return count > 0
}
