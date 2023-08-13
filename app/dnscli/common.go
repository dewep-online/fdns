package dnscli

import (
	"github.com/osspkg/fdns/app/db"
	"github.com/osspkg/goppy/plugins"
)

var Plugin = plugins.Plugin{
	Inject: func(v db.Connect) (Client, *object) {
		obj := NewClient(v)
		return obj, obj
	},
}
