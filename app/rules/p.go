package rules

import "go.osspkg.com/goppy/v2/plugins"

var Plugins = plugins.Plugins{}.Inject(
	NewRules,
	NewAdBlock,
	NewRegexpRules,
	NewStaticRules,
)
