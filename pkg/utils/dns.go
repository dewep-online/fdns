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
