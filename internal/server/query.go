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

package server

import (
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"

	"fdns/internal/utils"
)

func (s *Server) handler(writer dns.ResponseWriter, msg *dns.Msg) {

	response := &dns.Msg{}
	response.Authoritative = true
	response.SetReply(msg)
	response.SetRcode(msg, dns.RcodeSuccess)

	for _, q := range msg.Question {

		r := func() []dns.RR {
			switch q.Qtype {
			case dns.TypeA:
				return s.query(q.Name, s.serv.App.GetA, msg)

			case dns.TypeAAAA:
				return s.query(q.Name, s.serv.App.GetAAAA, msg)

			default:
				return s.query(q.Name, nil, msg)
			}
		}()
		response.Answer = append(response.Answer, r...)
	}

	if err := writer.WriteMsg(response); err != nil {
		logrus.WithError(err).WithField("dns", msg.String()).Error("response")
	}
}

func (s *Server) query(name string, fn func(string) (dns.RR, error), m *dns.Msg) (answer []dns.RR) {
	if fn != nil {
		rr, err := fn(name)
		switch err {
		case nil:
			answer = append(answer, rr)
			return
		case utils.ErrorEmptyIP:
			return
		}
	}

	nslookup := func(rm *dns.Msg) {
		for _, answ := range rm.Answer {

			switch answ.(type) {
			case *dns.A:
				if s.serv.App.InBlacklist(answ.(*dns.A).A) {
					if bha, er := s.serv.App.BlackHole(name); er != nil {
						answer = append(answer, bha)
					}
					continue
				}
				s.serv.App.Update(name, answ.(*dns.A).A.String(), "")
			case *dns.AAAA:
				if s.serv.App.InBlacklist(answ.(*dns.AAAA).AAAA) {
					continue
				}
				s.serv.App.Update(name, "", answ.(*dns.AAAA).AAAA.String())
			}

			answer = append(answer, answ)
		}
	}

	if ips, er := s.serv.App.GetDNS(name); er == nil {
		if rm, er := s.serv.Client.Exchange(m, ips); er == nil {
			nslookup(rm)
			return
		}
	}

	if rm, er := s.serv.Client.ExchangeRandomDNS(m); er == nil {
		nslookup(rm)
	}

	return
}
