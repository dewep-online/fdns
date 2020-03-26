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

	"github.com/sirupsen/logrus"

	"github.com/miekg/dns"

	"fdns/internal/app"
)

type DNSClient struct {
	ips    []string
	client *dns.Client
}

func New(c *app.ConfigApp) *DNSClient {
	rand.Seed(time.Now().UnixNano())

	return &DNSClient{
		client: &dns.Client{
			Net:          "tcp",
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
		},
		ips: c.DNS,
	}
}

func (d *DNSClient) Exchange(msg *dns.Msg) (*dns.Msg, error) {
	ns := d.ips[rand.Intn(len(d.ips))]

	resp, _, err := d.client.Exchange(msg, ns)

	var (
		mq []string
		mr []string
	)
	for _, q := range msg.Question {
		mq = append(mq, q.String())
	}

	if err != nil {
		logrus.WithError(err).
			WithField("ns", ns).
			WithField("q", strings.Join(mq, ",")).
			Info("reverse")
	} else {
		for _, a := range resp.Answer {
			mr = append(mr, a.String())
		}
		logrus.
			WithField("ns", ns).
			WithField("q", strings.Join(mq, ",")).
			WithField("r", strings.Join(mr, ",")).
			Info("reverse")
	}

	return resp, err
}

func (d *DNSClient) Up() error {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn("example.com."), dns.TypeNS)

	for _, ns := range d.ips {
		logrus.WithField("ns", ns).Info("check nslookup")
		if _, _, err := d.client.Exchange(msg, ns); err != nil {
			return err
		}
	}
	return nil
}
func (d *DNSClient) Down() error {
	return nil
}
