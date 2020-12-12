/*
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package dnscli

import (
	"math/rand"
	"strings"
	"time"

	"fdns/app/pkg/utils"

	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

type Client struct {
	ips []string
	cli *dns.Client
}

func New(c *Config) *Client {
	rand.Seed(time.Now().UnixNano())

	return &Client{
		cli: &dns.Client{
			Net:          "",
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
		},
		ips: c.DNS,
	}
}

func (o *Client) ExchangeRandomDNS(msg *dns.Msg) (*dns.Msg, error) {
	ns := o.ips[rand.Intn(len(o.ips))]

	return o.Exchange(msg, []string{ns})
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

		logrus.
			WithField("ns", ns).
			WithField("q", strings.Join(mq, ",")).
			WithField("r", strings.Join(mr, ",")).
			Info("reverse")

		break
	}

	if err != nil {
		logrus.WithError(err).
			WithField("ns", addrs).
			WithField("q", strings.Join(mq, ",")).
			Info("reverse")
	}

	return
}

func (o *Client) Up() error {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn("example.com."), dns.TypeNS)

	list := utils.ValidateIPs(o.ips)
	o.ips = make([]string, 0)

	for _, ip := range list {
		if _, _, err := o.cli.Exchange(msg, ip); err == nil {
			logrus.WithField("dns", ip).Info("add dns")
			o.ips = append(o.ips, ip)
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
