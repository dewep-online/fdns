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
	"net/http"
	"time"

	dns_client "fdns/internal/dns-client"

	"github.com/deweppro/core/pkg/app"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	app2 "fdns/internal/app"
)

type Server struct {
	serv *Services
	tcp  *dns.Server
	udp  *dns.Server
}

type Services struct {
	App    *app2.App
	Config *ConfigServer
	Client *dns_client.DNSClient
	Fc     *app.ForceClose
}

func New(s *Services) *Server {
	return &Server{
		serv: s,
	}
}

func (s *Server) Up() error {
	handler := dns.NewServeMux()
	handler.HandleFunc(".", s.handler)

	s.tcp = &dns.Server{Addr: s.serv.Config.Server.Addr,
		Net:          "tcp",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	s.udp = &dns.Server{Addr: s.serv.Config.Server.Addr,
		Net:          "udp",
		Handler:      handler,
		UDPSize:      65535,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		if err := s.tcp.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Error("tcp server")
			s.serv.Fc.Close()
		}
	}()
	go func() {
		if err := s.udp.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Error("udp server")
			s.serv.Fc.Close()
		}
	}()

	logrus.WithField("addr", s.serv.Config.Server.Addr).Info("dns serv start")
	return nil
}

func (s *Server) Down() (err error) {
	if err1 := s.tcp.Shutdown(); err1 != nil {
		err = err1
	}

	if err2 := s.udp.Shutdown(); err2 != nil {
		if err == nil {
			err = err2
		} else {
			err = errors.Wrap(err, err2.Error())
		}
	}

	logrus.WithField("addr", s.serv.Config.Server.Addr).Info("dns serv stop")
	return
}
