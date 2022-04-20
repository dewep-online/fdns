package blacklist

import (
	"net"
)

type Repository struct {
	conf           *Config
	blacklistIP    []net.IP
	blacklistIPNet []*net.IPNet
}

func New(c *Config) *Repository {
	return &Repository{
		conf: c,
	}
}

func (v *Repository) Up() error {
	for _, ip := range v.conf.BlackListIP {
		if _, n, err := net.ParseCIDR(ip); err == nil {
			v.blacklistIPNet = append(v.blacklistIPNet, n)
		} else {
			v.blacklistIP = append(v.blacklistIP, net.ParseIP(ip))
		}
	}
	return nil
}

func (v *Repository) Down() error {
	return nil
}

func (v *Repository) Has(ip net.IP) bool {
	for _, item := range v.blacklistIP {
		if item.Equal(ip) {
			return true
		}
	}
	for _, item := range v.blacklistIPNet {
		if item.Contains(ip) {
			return true
		}
	}
	return false
}
