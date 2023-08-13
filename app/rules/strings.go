package rules

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/osspkg/fdns/app/db"
	"github.com/osspkg/go-sdk/app"
	"github.com/osspkg/go-sdk/log"
	"github.com/osspkg/go-sdk/orm"
	"github.com/osspkg/go-sdk/routine"
)

type StringsBlock struct {
	db   db.Connect
	data []string
	mux  sync.RWMutex
}

func NewStringsBlock(dbc db.Connect) *StringsBlock {
	return &StringsBlock{
		db:   dbc,
		data: make([]string, 0, 100),
	}
}

func (v *StringsBlock) Up(ctx app.Context) error {
	routine.Interval(ctx.Context(), time.Hour, func(ctx context.Context) {
		if err := v.ForceUpdate(ctx); err != nil {
			log.WithError("err", err).Errorf("StringsBlock update")
		}
	})
	return nil
}

func (v *StringsBlock) Down() error {
	return nil
}

func (v *StringsBlock) ForceUpdate(ctx context.Context) error {
	result := make([]string, 0, 100)
	err := v.db.QueryContext("", ctx, func(q orm.Querier) {
		q.SQL("SELECT `data` FROM `strings_block`;")
		q.Bind(func(bind orm.Scanner) error {
			var data string
			if err := bind.Scan(&data); err != nil {
				return err
			}
			result = append(result, data)
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

func (v *StringsBlock) Contain(domain string) bool {
	v.mux.RLock()
	defer v.mux.RUnlock()

	for _, datum := range v.data {
		if strings.Contains(domain, datum) {
			return true
		}
	}
	return false
}
