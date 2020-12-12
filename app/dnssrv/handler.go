/*
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package dnssrv

import (
	"fdns/app/pkg/utils"

	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

func (o *Server) handler(writer dns.ResponseWriter, msg *dns.Msg) {
	response := &dns.Msg{}
	response.Authoritative = true
	response.SetReply(msg)
	response.SetRcode(msg, dns.RcodeSuccess)

	for _, q := range msg.Question {

		r := func() []dns.RR {
			switch q.Qtype {
			case dns.TypeA:
				return o.query(q.Name, o.store.GetA, msg)

			case dns.TypeAAAA:
				return o.query(q.Name, o.store.GetAAAA, msg)

			default:
				return o.query(q.Name, nil, msg)
			}
		}()
		response.Answer = append(response.Answer, r...)
	}

	if err := writer.WriteMsg(response); err != nil {
		logrus.WithError(err).WithField("dns", msg.String()).Error("response")
	}
}

func (o *Server) query(name string, fn func(string) ([]dns.RR, error), m *dns.Msg) (answer []dns.RR) {
	if fn != nil {
		rr, err := fn(name)
		switch err {
		case nil:
			answer = append(answer, rr...)
			return
		case utils.ErrEmptyIP:
			return
		}
	}

	if rm, err := o.store.DNS(name, m); err != nil {
		logrus.WithError(err).WithField("query", name).Error("nslookup")
	} else {
		answer = append(answer, rm...)
	}
	return
}
