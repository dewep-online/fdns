/*
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package utils

import (
	"net"

	"github.com/miekg/dns"
)

const (
	TTL = 3600
)

func CreateA(name string, ips []string) (result []dns.RR) {
	for _, ip := range ips {
		result = append(result, &dns.A{
			Hdr: dns.RR_Header{
				Name:   name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    uint32(TTL),
			},
			A: net.ParseIP(ip),
		})
	}
	return
}

func CreateAAAA(name string, ips []string) (result []dns.RR) {
	for _, ip := range ips {
		result = append(result, &dns.AAAA{
			Hdr: dns.RR_Header{
				Name:   name,
				Rrtype: dns.TypeAAAA,
				Class:  dns.ClassINET,
				Ttl:    uint32(TTL),
			},
			AAAA: net.ParseIP(ip),
		})
	}
	return
}
