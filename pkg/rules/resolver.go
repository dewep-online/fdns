package rules

import (
	"regexp"

	"github.com/dewep-online/fdns/pkg/utils"
)

type Resolver struct {
	reg    *regexp.Regexp
	format string
	types  uint
	ip4    []string
	ip6    []string
}

func (v *Resolver) Compile(name string) ([]string, []string, bool) {
	matches := v.reg.FindStringSubmatchIndex(name)
	if matches == nil {
		return nil, nil, false
	}
	if v.types == TypeDNS {
		return v.ip4, v.ip6, true
	}
	value := v.reg.ExpandString([]byte{}, v.format, name, matches)
	ip4, ip6 := utils.DecodeIPs(string(value))
	return ip4, ip6, true
}

func (v *Resolver) Match(name string) ([]string, []string, bool) {
	if !v.reg.MatchString(name) {
		return nil, nil, false
	}
	return v.ip4, v.ip6, true
}
