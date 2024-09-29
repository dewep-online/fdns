package rules

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/osspkg/fdns/app/db"
	"go.osspkg.com/goppy/v2/orm"
	"go.osspkg.com/ioutils/cache"
	"go.osspkg.com/logx"
	"go.osspkg.com/routine"
	"go.osspkg.com/xc"
)

type StaticRules struct {
	db   db.Connect
	data cache.TCacheReplace[string, []string]
}

func NewStaticRules(dbc db.Connect) *StaticRules {
	return &StaticRules{
		db:   dbc,
		data: cache.NewWithReplace[string, []string](),
	}
}

func (v *StaticRules) Up(ctx xc.Context) error {
	routine.Interval(ctx.Context(), time.Hour, func(ctx context.Context) {
		if err := v.ForceUpdate(ctx); err != nil {
			logx.Error("StaticRules update", "err", err)
		}
	})
	return nil
}

func (v *StaticRules) Down() error {
	return nil
}

func (v *StaticRules) key(qtype uint16, domain string) string {
	return fmt.Sprintf("%d %s", qtype, domain)
}

func (v *StaticRules) ForceUpdate(ctx context.Context) error {
	result := make(map[string][]string, 100)
	err := v.db.Main().Query(ctx, "load_static_rules", func(q orm.Querier) {
		q.SQL("SELECT `rule`,`qtype`,`data` FROM `static_regexp_rules` WHERE `deleted_at` IS NULL;")
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
			result[v.key(qtype, rule)] = data
			return nil
		})
	})
	if err != nil {
		return err
	}
	v.data.Replace(result)
	return nil
}

func (v *StaticRules) Convert(qtype uint16, domain string) []string {
	value, _ := v.data.Get(v.key(qtype, domain))
	return value
}
