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

func ValidateIP(ip string) (string, error) {
	if _, _, err := net.SplitHostPort(ip); err != nil {
		if v := net.ParseIP(ip); v != nil {
			return net.JoinHostPort(ip, "53"), nil
		}
		return "", ErrInvalidIP
	}
	return ip, nil
}

func ValidateIPs(list []string) (result []string) {
	for _, ip := range list {
		if v, er := ValidateIP(ip); er == nil {
			result = append(result, v)
		}
	}
	return
}

func ParseIPs(data string) (ip4, ip6 []string) {
	list := strings.Split(data, ",")
	for _, ip := range list {
		ip = strings.TrimSpace(ip)
		if _, err := ValidateIP(ip); err != nil {
			continue
		}
		if strings.Contains(ip, ":") {
			ip6 = append(ip6, ip)
			continue
		}
		ip4 = append(ip4, ip)
	}
	return
}
