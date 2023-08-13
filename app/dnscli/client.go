package dnscli

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/osspkg/fdns/app/ips"

	"github.com/osspkg/fdns/app/utils"

	"github.com/osspkg/go-sdk/errors"

	"github.com/osspkg/fdns/app/domain"

	"github.com/miekg/dns"
	"github.com/osspkg/fdns/app/db"
	"github.com/osspkg/go-sdk/app"
	"github.com/osspkg/go-sdk/log"
	"github.com/osspkg/go-sdk/orm"
	"github.com/osspkg/go-sdk/routine"
)

type (
	Rules struct {
		data map[string]map[string]struct{}
	}
	RulesGetter interface {
		Get(domain string) []string
	}
)

func NewRules() *Rules {
	return &Rules{data: make(map[string]map[string]struct{}, 100)}
}

func (v *Rules) Set(zone string, ips []string) {
	_, ok := v.data[zone]
	if !ok {
		v.data[zone] = make(map[string]struct{}, 2)
	}

	for _, ip := range ips {
		v.data[zone][ip] = struct{}{}
	}
}

func (v *Rules) Get(zone string) []string {
	for i := 2; i >= 0; i-- {
		vv := domain.Level(zone, i)
		if data, ok := v.data[vv]; ok {
			result := make([]string, 0, len(data))
			for ip := range data {
				result = append(result, ip)
			}
			return utils.Shuffle(result)
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var (
	ttlUpdate = 15 * time.Minute
)

type (
	Client interface {
		ForceUpdate(ctx context.Context) error
		Exchange(question dns.Question) ([]dns.RR, bool)
		ExchangeCustom(question dns.Question, adders []string) ([]dns.RR, bool)
	}
	object struct {
		rules RulesGetter
		cli   *dns.Client
		db    db.Connect
		mux   sync.RWMutex
	}
)

func NewClient(dbc db.Connect) *object {
	return &object{
		cli: &dns.Client{
			Net:          "",
			ReadTimeout:  time.Second * 3,
			WriteTimeout: time.Second * 3,
		},
		rules: NewRules(),
		db:    dbc,
	}
}

func (v *object) Up(ctx app.Context) error {
	routine.Interval(ctx.Context(), ttlUpdate, func(ctx context.Context) {
		if err := v.ForceUpdate(ctx); err != nil {
			log.WithError("err", err).Errorf("DNS Client update dns list")
		}
	})
	return nil
}

func (v *object) Down() error {
	return nil
}

func (v *object) ForceUpdate(ctx context.Context) error {
	result := NewRules()
	err := v.db.QueryContext("", ctx, func(q orm.Querier) {
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
			result.Set(zone, ips.NormalizeDNS(data...))
			return nil
		})
	})
	if err != nil {
		return err
	}
	v.mux.Lock()
	v.rules = result
	v.mux.Unlock()
	return nil
}

func (v *object) Exchange(question dns.Question) ([]dns.RR, bool) {
	return v.ExchangeCustom(question, v.rules.Get(question.Name))
}

func (v *object) ExchangeCustom(question dns.Question, adders []string) ([]dns.RR, bool) {
	var (
		mr   []string
		errs error
	)

	msg := new(dns.Msg).SetQuestion(question.Name, question.Qtype)

	for _, ns := range adders {
		resp, _, err := v.cli.Exchange(msg, ns)
		if err != nil {
			errs = errors.Wrap(errs, err)
			continue
		}

		for _, a := range resp.Answer {
			mr = append(mr, a.String())

		}

		log.WithFields(log.Fields{
			"dns":      ns,
			"question": question.String(),
			"answer":   strings.Join(mr, ","),
		}).Infof("DNS Client exchange")

		return resp.Answer, len(resp.Answer) > 0
	}

	if errs != nil {
		log.WithFields(log.Fields{
			"dns":      strings.Join(adders, ","),
			"question": question.String(),
			"err":      errs.Error(),
		}).Errorf("DNS Client exchange")
	}

	return nil, false
}
