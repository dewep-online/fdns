package dnsserver

import (
	"github.com/dewep-online/fdns/pkg/utils"
	"github.com/deweppro/go-logger"
	"github.com/miekg/dns"
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
		logger.Errorf("response ERROR: %s DNS %s", err.Error(), msg.String())
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
		logger.Errorf("nslookup ERROR: %s QUERY %s", err.Error(), name)
	} else {
		answer = append(answer, rm...)
	}
	return
}
