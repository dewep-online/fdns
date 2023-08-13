package dns

import (
	"net/http"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/osspkg/fdns/app/resolver"
	"github.com/osspkg/go-sdk/app"
	"github.com/osspkg/go-sdk/errors"
	"github.com/osspkg/go-sdk/log"
)

type Server struct {
	conf ConfigItem
	serv []*dns.Server
	res  resolver.Resolver
	wg   sync.WaitGroup
}

func NewServer(conf ConfigItem, res resolver.Resolver) *Server {
	return &Server{
		conf: conf,
		serv: make([]*dns.Server, 0, 2),
		res:  res,
	}
}

func (v *Server) Up(ctx app.Context) error {
	handler := dns.NewServeMux()
	handler.HandleFunc(".", v.DNSHandler)

	v.serv = append(v.serv, &dns.Server{
		Addr:         v.conf.Addr,
		Net:          "tcp",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
	v.serv = append(v.serv, &dns.Server{
		Addr:         v.conf.Addr,
		Net:          "udp",
		Handler:      handler,
		UDPSize:      65535,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	for _, s := range v.serv {
		s := s
		v.wg.Add(1)

		log.WithFields(log.Fields{
			"address": s.Addr,
			"net":     s.Net,
		}).Infof("Run DNS Server")

		go func(srv *dns.Server) {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.WithFields(log.Fields{
					"err":     err.Error(),
					"address": srv.Addr,
					"net":     srv.Net,
				}).Errorf("Run DNS Server")
				ctx.Close()
			}
			v.wg.Done()
		}(s)
	}

	return nil
}

func (v *Server) Down() error {
	for _, s := range v.serv {
		if err := s.Shutdown(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithFields(log.Fields{
				"err":     err.Error(),
				"address": s.Addr,
				"net":     s.Net,
			}).Errorf("Shutdown DNS Server")
			continue
		}
		log.WithFields(log.Fields{
			"address": s.Addr,
			"net":     s.Net,
		}).Infof("Shutdown DNS Server")
	}

	v.wg.Wait()
	return nil
}

func (v *Server) DNSHandler(w dns.ResponseWriter, msg *dns.Msg) {
	response := &dns.Msg{}
	response.Authoritative = true
	response.RecursionAvailable = true
	response.SetReply(msg)
	response.SetRcode(msg, dns.RcodeSuccess)

	for _, question := range msg.Question {
		if question.Qtype != dns.TypeA {
			continue
		}
		response.Answer = append(response.Answer, v.res.Resolve(question)...)
	}

	if err := w.WriteMsg(response); err != nil {
		log.WithFields(log.Fields{
			"err":      err.Error(),
			"question": msg.String(),
			"answer":   response.String(),
		}).Errorf("DNS handler", err.Error(), msg.String())
	}
}
