package dnsserver

import (
	"github.com/deweppro/go-logger"
	"github.com/miekg/dns"
)

func (v *Server) handler(writer dns.ResponseWriter, msg *dns.Msg) {
	response := &dns.Msg{}
	response.Authoritative = true
	response.RecursionAvailable = true
	response.SetReply(msg)
	response.SetRcode(msg, dns.RcodeSuccess)

	if len(msg.Question) > 0 {
		response.Answer = append(response.Answer, v.store.Resolve(msg.Question[0])...)
	}

	if err := writer.WriteMsg(response); err != nil {
		logger.Errorf("response ERROR: %s DNS %s", err.Error(), msg.String())
	}
}
