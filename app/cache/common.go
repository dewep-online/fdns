package cache

import (
	"github.com/osspkg/goppy/plugins"
)

var Plugins = plugins.Plugins{}.Inject(
	NewRecords,
)
