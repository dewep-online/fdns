package dnscli

import (
	"context"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/dewep-online/fdns/pkg/utils"
	"github.com/deweppro/go-app/application/ctx"

	"github.com/dewep-online/fdns/pkg/database"

	"github.com/deweppro/go-logger"
	"github.com/miekg/dns"
)

type Client struct {
	ips []string
	cli *dns.Client
	db  *database.Database
	mux sync.RWMutex
}

func New(db *database.Database) *Client {
	return &Client{
		cli: &dns.Client{
			Net:          "",
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
		},
		ips: make([]string, 0, 10),
		db:  db,
	}
}

func (o *Client) Up(ctx ctx.Context) error {
	utils.Interval(ctx.Context(), time.Minute*15, o.UpgradeNS)
	return nil
}

func (o *Client) Down(_ ctx.Context) error {
	return nil
}

func (o *Client) UpgradeNS(ctx context.Context) {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn("example.com."), dns.TypeNS)

	rules, err := o.db.GetRules(ctx, database.NS, 0, database.ActiveTrue)
	if err != nil {
		logger.Errorf("load dns: %s", err.Error())
		return
	}

	var ips []string
	for _, rule := range rules {
		ip4, ip6 := utils.DecodeIPs(rule.IPs)
		ips = append(ips, ip4...)
		ips = append(ips, ip6...)
	}

	ips = utils.ValidateDNSs(ips)
	result := make([]string, 0, len(ips))

	for _, ip := range ips {
		if _, _, err = o.cli.Exchange(msg, ip); err != nil {
			logger.Errorf("error dns: [%s] %s", ip, err.Error())
			continue
		}
		logger.Infof("set dns: %s", ip)
		result = append(result, ip)
	}

	o.mux.Lock()
	o.ips = result
	o.mux.Unlock()
}

func (o *Client) ExchangeRandomDNS(msg *dns.Msg) (*dns.Msg, error) {
	o.mux.RLock()
	defer o.mux.RUnlock()

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(o.ips), func(i, j int) { o.ips[i], o.ips[j] = o.ips[j], o.ips[i] })

	return o.Exchange(msg, o.ips)
}

func (o *Client) Exchange(msg *dns.Msg, addrs []string) (resp *dns.Msg, err error) {
	var mq, mr []string

	for _, ns := range addrs {
		resp, _, err = o.cli.Exchange(msg, ns)
		if err != nil {
			continue
		}

		for _, q := range msg.Question {
			mq = append(mq, q.String())
		}

		for _, a := range resp.Answer {
			mr = append(mr, a.String())
		}

		logger.Infof("reverse: NS: %s QUERY: %s RESPONSE: %s",
			ns, strings.Join(mq, ","), strings.Join(mr, ","))

		break
	}

	if err != nil {
		logger.Infof("reverse: NS: %s QUERY: %s ERROR: %s",
			addrs, strings.Join(mq, ","), err.Error())
	}

	return
}
