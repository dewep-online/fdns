package dnsserver

import (
	"net/http"
	"time"

	"github.com/dewep-online/fdns/pkg/rules"
	"github.com/deweppro/go-app/application"
	"github.com/deweppro/go-logger"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

type (
	Server struct {
		tcp   *dns.Server
		udp   *dns.Server
		conf  *ConfigTCP
		close *application.ForceClose
		store *rules.Repository
	}
)

func New(c *ConfigTCP, f *application.ForceClose, r *rules.Repository) *Server {
	return &Server{
		conf:  c,
		close: f,
		store: r,
	}
}

func (o *Server) Up() error {
	handler := dns.NewServeMux()
	handler.HandleFunc(".", o.handler)

	o.tcp = &dns.Server{Addr: o.conf.Server.Addr,
		Net:          "tcp",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	o.udp = &dns.Server{Addr: o.conf.Server.Addr,
		Net:          "udp",
		Handler:      handler,
		UDPSize:      65535,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		if err := o.tcp.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("dns tcp server: %s", err.Error())
			o.close.Close()
		}
	}()
	go func() {
		if err := o.udp.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("dns udp server: %s", err.Error())
			o.close.Close()
		}
	}()

	logger.Infof("dns serv start: %s", o.conf.Server.Addr)
	return nil
}

func (o *Server) Down() (err error) {
	if err1 := o.tcp.Shutdown(); err1 != nil {
		err = err1
	}

	if err2 := o.udp.Shutdown(); err2 != nil {
		if err == nil {
			err = err2
		} else {
			err = errors.Wrap(err, err2.Error())
		}
	}

	logger.Infof("dns serv stop: %s", o.conf.Server.Addr)
	return
}
