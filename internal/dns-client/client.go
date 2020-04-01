/*
 * Copyright (c) 2020.  Mikhail Knyazhev <markus621@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/gpl-3.0.html>.
 */

package dns_client

import (
	"math/rand"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"

	"fdns/internal/app"
	"fdns/internal/utils"
)

type DNSClient struct {
	ips    []string
	client *dns.Client
}

func New(c *app.ConfigApp) *DNSClient {
	rand.Seed(time.Now().UnixNano())

	return &DNSClient{
		client: &dns.Client{
			Net:          "",
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
		},
		ips: c.DNS,
	}
}

func (d *DNSClient) ExchangeRandomDNS(msg *dns.Msg) (*dns.Msg, error) {
	ns := d.ips[rand.Intn(len(d.ips))]

	return d.Exchange(msg, []string{ns})
}

func (d *DNSClient) Exchange(msg *dns.Msg, addrs []string) (resp *dns.Msg, err error) {
	var mq, mr []string

	for _, ns := range addrs {
		resp, _, err = d.client.Exchange(msg, ns)
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

func (d *DNSClient) Up() error {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn("example.com."), dns.TypeNS)

	list := utils.ValidateIPs(d.ips)
	d.ips = make([]string, 0)

	for _, ip := range list {
		if _, _, err := d.client.Exchange(msg, ip); err == nil {
			logrus.WithField("dns", ip).Info("add dns")
			d.ips = append(d.ips, ip)
		}
	}

	if len(d.ips) == 0 {
		return utils.ErrorEmptyDNSList
	}
	return nil
}

func (d *DNSClient) Down() error {
	return nil
}
