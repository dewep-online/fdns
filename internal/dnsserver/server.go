package dnsserver

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dewep-online/fdns/pkg/rules"
	"github.com/deweppro/go-app/application/ctx"
	"github.com/deweppro/go-logger"
	"github.com/miekg/dns"
)

type (
	Server struct {
		conf    *ConfigTCP
		servers map[string]*dns.Server
		store   *rules.Repository
	}
)

func New(c *ConfigTCP, r *rules.Repository) *Server {
	return &Server{
		servers: make(map[string]*dns.Server),
		conf:    c,
		store:   r,
	}
}

func (v *Server) Up(cx ctx.Context) error {
	handler := dns.NewServeMux()
	handler.HandleFunc(".", v.handler)

	if err := v.dnsSetup(handler); err != nil {
		return err
	}

	if err := v.dotSetup(handler); err != nil {
		return err
	}

	v.runServers(cx)

	return nil

}

func (v *Server) Down(_ ctx.Context) (err error) {
	for name, server := range v.servers {
		if err := server.Shutdown(); err != nil && err != http.ErrServerClosed {
			logger.WithFields(logger.Fields{
				"err":  err.Error(),
				"name": name,
				"ip":   server.Addr,
			}).Errorf("shutdown server")
		} else {
			logger.WithFields(logger.Fields{
				"name": name,
				"ip":   server.Addr,
			}).Infof("shutdown server")
		}

	}
	return nil
}

func (v *Server) dnsSetup(handler *dns.ServeMux) error {
	if !v.conf.Srv.Enable {
		return nil
	}

	v.servers["dns_tcp"] = &dns.Server{
		Addr:         v.conf.Srv.Addr,
		Net:          "tcp",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	v.servers["dns_udp"] = &dns.Server{
		Addr:         v.conf.Srv.Addr,
		Net:          "udp",
		Handler:      handler,
		UDPSize:      65535,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return nil
}

func (v *Server) dotSetup(handler *dns.ServeMux) error {
	if !v.conf.DoT.Enable {
		return nil
	}

	cert, err := v.tlsCertificate()
	if err != nil {
		return fmt.Errorf("read tls certificate: %w", err)
	}

	certConf := &tls.Config{Certificates: []tls.Certificate{*cert}}

	v.servers["dot_tcp"] = &dns.Server{
		Addr:         v.conf.DoT.Addr,
		Net:          "tcp",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		TLSConfig:    certConf,
	}

	v.servers["dot_udp"] = &dns.Server{
		Addr:         v.conf.DoT.Addr,
		Net:          "udp",
		Handler:      handler,
		UDPSize:      65535,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		TLSConfig:    certConf,
	}

	return nil
}

func (v *Server) tlsCertificate() (*tls.Certificate, error) {
	pubPem, err := ioutil.ReadFile(v.conf.DoT.Cert.Public)
	if err != nil {
		return nil, err
	}

	keyPem, err := ioutil.ReadFile(v.conf.DoT.Cert.Private)
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(pubPem, keyPem)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

func (v *Server) runServers(cx ctx.Context) {
	for name, server := range v.servers {
		server := server
		name := name

		logger.WithFields(logger.Fields{
			"name": name,
			"ip":   server.Addr,
		}).Infof("start server")

		go func(name string, srv *dns.Server) {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.WithFields(logger.Fields{
					"err":  err.Error(),
					"name": name,
					"ip":   srv.Addr,
				}).Errorf("start server")
				cx.Close()
			}
		}(name, server)
	}
}
