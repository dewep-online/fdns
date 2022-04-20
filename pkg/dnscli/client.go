package dnscli

import (
	"math/rand"
	"strings"
	"time"

	"github.com/dewep-online/fdns/pkg/utils"
	"github.com/deweppro/go-logger"
	"github.com/miekg/dns"
)

type Client struct {
	ips []string
	cli *dns.Client
}

func New(c *Config) *Client {
	return &Client{
		cli: &dns.Client{
			Net:          "",
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
		},
		ips: append([]string{}, c.DNS...),
	}
}

func (o *Client) ExchangeRandomDNS(msg *dns.Msg) (*dns.Msg, error) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(o.ips), func(i, j int) { o.ips[i], o.ips[j] = o.ips[j], o.ips[i] })

	return o.Exchange(msg, o.ips)
}

func (o *Client) Exchange(msg *dns.Msg, addrs []string) (*dns.Msg, error) {
	var (
		mq, mr []string
		resp   *dns.Msg
		err    error
	)

	for _, ns := range addrs {

		if resp, _, err = o.cli.Exchange(msg, ns); err != nil {

			logger.WithFields(logger.Fields{
				"ns":  ns,
				"q":   strings.Join(mq, ","),
				"err": err.Error(),
			}).Errorf("receive ip")

			continue
		}

		for _, q := range msg.Question {
			mq = append(mq, q.String())
		}

		for _, a := range resp.Answer {
			mr = append(mr, a.String())
		}

		logger.WithFields(logger.Fields{
			"ns": ns,
			"q":  strings.Join(mq, ","),
			"a":  strings.Join(mr, ","),
		}).Infof("receive ip")

		break
	}

	if resp == nil && err == nil {
		return nil, utils.ErrEmptyIP
	}

	return resp, nil
}

func (o *Client) Up() error {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn("example.com."), dns.TypeNS)

	list := utils.ValidateDNSs(o.ips)
	o.ips = make([]string, 0)

	for _, ip := range list {
		logger.WithFields(logger.Fields{
			"ip": ip,
		}).Infof("add ns")
		o.ips = append(o.ips, ip)

		if _, _, err := o.cli.Exchange(msg, ip); err != nil {
			logger.WithFields(logger.Fields{
				"err": err.Error(),
				"ip":  ip,
			}).Errorf("add ns")
		}
	}

	if len(o.ips) == 0 {
		return utils.ErrEmptyDNSList
	}
	return nil
}

func (o *Client) Down() error {
	return nil
}
