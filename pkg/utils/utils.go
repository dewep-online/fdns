package utils

import (
	"errors"
	"net"
	"strings"
)

var (
	ErrInvalidIP     = errors.New("invalid ip")
	ErrEmptyDNSList  = errors.New("dns list is empty")
	ErrEmptyIP       = errors.New("ip is empty")
	ErrCacheNotFound = errors.New("cache is not found")
)

func ValidateDNS(ip string) (string, error) {
	if _, _, err := net.SplitHostPort(ip); err != nil {
		if v := net.ParseIP(ip); v != nil {
			return net.JoinHostPort(ip, "53"), nil
		}
		return "", ErrInvalidIP
	}
	return ip, nil
}

func ValidateDNSs(list []string) (result []string) {
	for _, ip := range list {
		if v, er := ValidateDNS(ip); er == nil {
			result = append(result, v)
		}
	}
	return
}

func DecodeIPs(data string) (ip4, ip6 []string) {
	list := strings.Split(data, ",")
	for _, host := range list {
		host = strings.TrimSpace(host)
		ip, port, err := net.SplitHostPort(host)
		if err != nil {
			ip = host
		}
		v := net.ParseIP(ip)
		if v == nil {
			continue
		}
		ip = v.String()
		if len(port) > 0 {
			host = net.JoinHostPort(ip, port)
		}
		if strings.Contains(ip, ":") {
			ip6 = append(ip6, host)
			continue
		}
		ip4 = append(ip4, host)
	}
	return
}

func EncodeIPs(ip4, ip6 []string) string {
	return strings.Join(append(ip4, ip6...), ", ")
}
