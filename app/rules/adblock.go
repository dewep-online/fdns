package rules

import (
	"context"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/osspkg/fdns/app/db"
	"github.com/osspkg/fdns/app/httpcli"
	"github.com/osspkg/fdns/app/utils"
	"github.com/osspkg/go-sdk/app"
	"github.com/osspkg/go-sdk/log"
	"github.com/osspkg/go-sdk/orm"
	"github.com/osspkg/go-sdk/routine"
)

var rex = regexp.MustCompile(`(?miU)^\|\|([a-z0-9-.]+)\^(\n|\r|\$)`)

type AdBlock struct {
	cli  *httpcli.Client
	db   db.Connect
	data map[string]struct{}
	mux  sync.RWMutex
}

func NewAdBlock(dbc db.Connect, hc *httpcli.Client) *AdBlock {
	return &AdBlock{
		cli:  hc,
		db:   dbc,
		data: make(map[string]struct{}, 100000),
	}
}

func (v *AdBlock) Up(ctx app.Context) error {
	if err := v.ForceUpdate(ctx.Context()); err != nil {
		log.WithError("err", err).Errorf("AdBlock update")
	}
	go routine.Interval(ctx.Context(), time.Hour, func(ctx context.Context) {
		if err := v.UpgradeRules(ctx); err != nil {
			log.WithError("err", err).Errorf("AdBlock upgrade")
		}
		if err := v.ForceUpdate(ctx); err != nil {
			log.WithError("err", err).Errorf("AdBlock update")
		}
	})
	return nil
}

func (v *AdBlock) Down() error {
	return nil
}

func (v *AdBlock) UpgradeRules(ctx context.Context) error {
	list := make(map[uint64]string, 100000)
	err := v.db.QueryContext("", ctx, func(q orm.Querier) {
		q.SQL("SELECT `id`,`data` FROM `adblock_list`;")
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
		_, b, err0 := v.cli.Get(uri)
		if err0 != nil {
			return err0
		}
		rexResult := rex.FindAll(b, -1)
		result := make([]string, 0, 100)
		for _, domain := range rexResult {
			rule := strings.Trim(string(domain[2:len(domain)-1]), "\n^") + "."
			if v.hash(rule) {
				continue
			}
			result = append(result, rule)
			if len(result) == 100 {
				if err0 = v.save(ctx, id, result); err0 != nil {
					return err0
				}
				result = result[:0]
			}
		}

		if len(result) > 0 {
			if err0 = v.save(ctx, id, result); err0 != nil {
				return err0
			}
		}

		log.WithFields(log.Fields{
			"uri":   uri,
			"count": len(rexResult),
		}).Infof("AdBlock upgrade")
	}

	return nil
}

func (v *AdBlock) hash(rule string) bool {
	v.mux.RLock()
	defer v.mux.RUnlock()
	_, ok := v.data[rule]
	return ok
}

func (v *AdBlock) save(ctx context.Context, id uint64, data []string) error {
	return v.db.TransactionContext("", ctx, func(v orm.Tx) {
		v.Exec(func(e orm.Executor) {
			e.SQL("INSERT IGNORE INTO `adblock_rules` (`list_id`, `data`, `hash`, `updated_at`)" +
				"VALUES (?, ?, ?, now());")
			for _, datum := range data {
				e.Params(id, datum, utils.Sha1(datum))
			}
		})
	})
}

func (v *AdBlock) ForceUpdate(ctx context.Context) error {
	result := make(map[string]struct{}, 100000)
	err := v.db.QueryContext("", ctx, func(q orm.Querier) {
		q.SQL("SELECT `data` FROM `adblock_rules` WHERE `deleted_at` IS NULL;")
		q.Bind(func(bind orm.Scanner) error {
			var data string
			if err := bind.Scan(&data); err != nil {
				return err
			}
			result[data] = struct{}{}
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

func (v *AdBlock) Contain(domain string) bool {
	v.mux.RLock()
	defer v.mux.RUnlock()

	for datum := range v.data {
		if strings.HasSuffix(domain, datum) {
			return true
		}
	}
	return false
}
