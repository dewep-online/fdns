package rules

import (
	"regexp"

	"github.com/dewep-games/fdns/pkg/utils"
)

type (
	RuleMap map[string]*Rule
	Rule    struct {
		reg    *regexp.Regexp
		format string
		ip4    []string
		ip6    []string
	}
)

func (o *Rule) Compile(name string) (ip4, ip6 []string, ok bool) {
	matches := o.reg.FindStringSubmatchIndex(name)
	if matches == nil {
		ok = false
		return
	}
	value := o.reg.ExpandString([]byte{}, o.format, name, matches)
	ip4, ip6 = utils.ParseIPs(string(value))
	ok = true
	return
}

func (o *Rule) Match(name string) ([]string, []string, bool) {
	if !o.reg.MatchString(name) {
		return nil, nil, false
	}
	return o.ip4, o.ip6, true
}
