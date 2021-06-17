package blacklist

import (
	"net"

	"github.com/dewep-online/fdns/pkg/utils"
	"github.com/miekg/dns"
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

func (o *Repository) Up() error {
	for _, ip := range o.conf.BlackListIP {
		if _, n, err := net.ParseCIDR(ip); err == nil {
			o.blacklistIPNet = append(o.blacklistIPNet, n)
		} else {
			o.blacklistIP = append(o.blacklistIP, net.ParseIP(ip))
		}
	}
	return nil
}

func (o *Repository) Down() error {
	return nil
}

func (o *Repository) Has(ip net.IP) bool {
	for _, item := range o.blacklistIP {
		if item.Equal(ip) {
			return true
		}
	}
	for _, item := range o.blacklistIPNet {
		if item.Contains(ip) {
			return true
		}
	}
	return false
}

func (o *Repository) BlackHole(name string) ([]dns.RR, error) {
	if len(o.conf.BlackHoleIP) == 0 {
		return nil, utils.ErrEmptyIP
	}
	return utils.CreateA(name, []string{o.conf.BlackHoleIP}), nil
}
