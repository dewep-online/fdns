package rules

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/deweppro/go-logger"

	"github.com/dewep-online/fdns/pkg/httpcli"

	"github.com/dewep-online/fdns/pkg/utils"
)

const (
	TypeNone   uint = 0
	TypeDNS    uint = 1
	TypeRegexp uint = 2
	TypeHost   uint = 3
)

type (
	HostSetter interface {
		SetHostResolve(domain string, ip4, ip6 []string, ttl int64)
	}
	ResolveSetter interface {
		SetRexResolve(format string, rx *regexp.Regexp, ip4, ip6 []string, tp uint)
	}
)

func HostRules(data map[string]string, setter HostSetter) error {
	for domain, ips := range data {
		ip4, ip6 := utils.DecodeIPs(ips)
		setter.SetHostResolve(domain+".", ip4, ip6, 0)
	}
	return nil
}

func DNSRules(data map[string]string, setter ResolveSetter) error {
	for domain, ips := range data {
		ip4, ip6 := utils.DecodeIPs(ips)

		domain = regexp.QuoteMeta(domain)
		domain = strings.ReplaceAll(domain, "\\?", "?")
		domain = strings.ReplaceAll(domain, "\\*", ".*")
		domain = fmt.Sprintf("^.*%s\\.$", strings.Trim(domain, "^$"))
		rx, err := regexp.Compile(domain)
		if err != nil {
			return err
		}

		setter.SetRexResolve(
			"",
			rx,
			utils.ValidateDNSs(ip4),
			utils.ValidateDNSs(ip6),
			TypeDNS,
		)
	}
	return nil
}

func RegexpRules(data map[string]string, setter ResolveSetter) error {
	for domain, ips := range data {
		ip4, ip6 := utils.DecodeIPs(ips)
		domain = fmt.Sprintf("^%s\\.$", strings.Trim(domain, "^$"))
		rx, err := regexp.Compile(domain)
		if err != nil {
			return err
		}

		setter.SetRexResolve(
			ips,
			rx,
			utils.ValidateDNSs(ip4),
			utils.ValidateDNSs(ip6),
			TypeRegexp,
		)
	}
	return nil
}

func QueryRules(data map[string]string, setter ResolveSetter) error {
	for domain, ips := range data {
		ip4, ip6 := utils.DecodeIPs(ips)

		domain = regexp.QuoteMeta(domain)
		domain = strings.ReplaceAll(domain, "\\?", ".")
		domain = strings.ReplaceAll(domain, "\\*", ".*")
		domain = fmt.Sprintf("^%s\\.$", strings.Trim(domain, "^$"))
		rx, err := regexp.Compile(domain)
		if err != nil {
			return err
		}

		setter.SetRexResolve(
			ips,
			rx,
			utils.ValidateDNSs(ip4),
			utils.ValidateDNSs(ip6),
			TypeRegexp,
		)
	}
	return nil
}

var (
	cli = httpcli.New()
	rex = regexp.MustCompile(`\|\|([a-z0-9-.]+)\^`)
)

func AdblockRules(data []string, setter HostSetter) {

	for _, uri := range data {
		code, b, err := cli.Call(http.MethodGet, uri, nil)
		if err != nil {
			logger.Warnf("adblock [%d] %s: %s", code, uri, err.Error())
			continue
		}
		if code != http.StatusOK {
			logger.Warnf("adblock [%d] %s: %s", code, uri, err.Error())
			continue
		}

		result := rex.FindAll(b, -1)
		logger.Infof("adblock [%d] %s", len(result), uri)
		for _, domain := range result {
			setter.SetHostResolve(string(domain[2:len(domain)-1])+".", nil, nil, 0)
		}
	}
}
