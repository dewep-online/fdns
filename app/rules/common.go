package rules

import (
	"github.com/osspkg/goppy/plugins"
)

var Plugins = plugins.Plugins{}.Inject(
	NewRules,
	NewStringsBlock,
	NewAdBlock,
)

type Rules struct {
	strings *StringsBlock
	adblock *AdBlock
}

func NewRules(sb *StringsBlock, ab *AdBlock) *Rules {
	return &Rules{
		strings: sb,
		adblock: ab,
	}
}

func (v *Rules) IsBlocked(domain string) bool {
	if v.strings.Contain(domain) {
		return true
	}
	if v.adblock.Contain(domain) {
		return true
	}
	return false
}
