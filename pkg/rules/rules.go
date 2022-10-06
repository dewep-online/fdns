package rules

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/dewep-online/fdns/pkg/httpcli"
	"github.com/dewep-online/fdns/pkg/utils"
	"github.com/deweppro/go-logger"
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
		SetRexResolve(rule, format string, rx *regexp.Regexp, ip4, ip6 []string, tp uint)
	}
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func HostRules(data map[string]string, setter HostSetter) error {
	for domain, ips := range data {
		ip4, ip6 := utils.DecodeIPs(ips)
		setter.SetHostResolve(domain, ip4, ip6, 0)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func DNSRules(data map[string]string, setter ResolveSetter) error {
	for rule, ips := range data {
		ip4, ip6 := utils.DecodeIPs(ips)

		domain := regexp.QuoteMeta(rule)
		domain = strings.ReplaceAll(domain, "\\?", "?")
		domain = strings.ReplaceAll(domain, "\\*", ".*")
		domain = fmt.Sprintf("^.*%s\\.$", strings.Trim(domain, "^$"))
		rx, err := regexp.Compile(domain)
		if err != nil {
			return err
		}

		setter.SetRexResolve(
			rule,
			"",
			rx,
			utils.ValidateDNSs(ip4),
			utils.ValidateDNSs(ip6),
			TypeDNS,
		)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func RegexpRules(data map[string]string, setter ResolveSetter) error {
	for rule, ips := range data {
		ip4, ip6 := utils.DecodeIPs(ips)
		domain := fmt.Sprintf("^%s\\.$", strings.Trim(rule, "^$"))
		rx, err := regexp.Compile(domain)
		if err != nil {
			return err
		}

		setter.SetRexResolve(
			rule,
			ips,
			rx,
			utils.ValidateDNSs(ip4),
			utils.ValidateDNSs(ip6),
			TypeRegexp,
		)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func QueryRules(data map[string]string, setter ResolveSetter) error {
	for rule, ips := range data {
		ip4, ip6 := utils.DecodeIPs(ips)

		domain := regexp.QuoteMeta(rule)
		domain = strings.ReplaceAll(domain, "\\?", ".")
		domain = strings.ReplaceAll(domain, "\\*", ".*")
		domain = fmt.Sprintf("^%s\\.$", strings.Trim(domain, "^$"))
		rx, err := regexp.Compile(domain)
		if err != nil {
			return err
		}

		setter.SetRexResolve(
			rule,
			ips,
			rx,
			utils.ValidateDNSs(ip4),
			utils.ValidateDNSs(ip6),
			TypeRegexp,
		)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var (
	cli = httpcli.New()
	rex = regexp.MustCompile(`\|\|([a-z0-9-.]+)\^(\n|\r)`)
)

func LoadAdblockRules(uri string) []string {
	utils.Retry(5, time.Second*5, func() error {
		u, err := url.Parse(uri)
		if err != nil {
			logger.Warnf("adblock-rules parse base url [%s]: %s", uri, utils.StringError(err))
			return err
		}
		_, err = net.LookupIP(u.Host)
		if err != nil {
			logger.Warnf("adblock-rules nslookup [%s]: %s", u.Host, utils.StringError(err))
		}
		return err
	})
	code, b, err := cli.Call(http.MethodGet, uri, nil)
	if err != nil || code != http.StatusOK {
		logger.Warnf("adblock-rules [%d] %s: %s", code, uri, utils.StringError(err))
		return nil
	}
	result := make([]string, 0, 10)
	if err = json.Unmarshal(b, &result); err != nil {
		logger.Warnf("adblock-rules [%d] %s: %s", code, uri, utils.StringError(err))
	}
	return result
}

func AdblockRules(data []string, setter func(uri string, domains []string)) {
	for _, uri := range data {
		code, b, err := cli.Call(http.MethodGet, uri, nil)
		if err != nil || code != http.StatusOK {
			logger.Warnf("adblock [%d] %s: %s", code, uri, utils.StringError(err))
			continue
		}

		rexResult := rex.FindAll(b, -1)
		result := make([]string, 0, len(rexResult))
		for _, domain := range rexResult {
			result = append(result,
				strings.Trim(string(domain[2:len(domain)-1]), "\n^")+".")
		}

		logger.Infof("adblock [%d] %s", len(result), uri)
		setter(uri, result)
	}
}
