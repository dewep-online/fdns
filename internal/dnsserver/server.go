package dnsserver

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/deweppro/go-app/application/ctx"

	"github.com/dewep-online/fdns/pkg/rules"
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

func (v *Server) Up(ctx ctx.Context) error {
	handler := dns.NewServeMux()
	handler.HandleFunc(".", v.handler)

	if err := v.dnsSetup(handler); err != nil {
		return err
	}

	if err := v.dotSetup(handler); err != nil {
		return err
	}

	v.runServers(ctx)

	return nil

}

func (v *Server) Down(_ ctx.Context) (err error) {
	for name, server := range v.servers {
		if err := server.Shutdown(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("%s [%s]: %s", name, server.Addr, err.Error())
		} else {
			logger.Infof("%s stop: %s", name, server.Addr)
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

func (v *Server) runServers(ctx ctx.Context) {
	for name, server := range v.servers {
		server := server
		name := name

		go func(name string, srv *dns.Server) {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Errorf("%s [%s]: %s", name, srv.Addr, err.Error())
				ctx.Close()
			}
		}(name, server)

		logger.Infof("%s start: %s", name, server.Addr)
	}
}
