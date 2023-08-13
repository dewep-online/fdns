package db

import (
	"github.com/osspkg/go-sdk/orm"
	"github.com/osspkg/goppy/plugins"
	"github.com/osspkg/goppy/plugins/database"
)

var Plugin = plugins.Plugin{
	Inject: func(v database.MySQL) Connect {
		return v.Pool("main")
	},
}

type (
	Connect interface {
		orm.Stmt
	}
)
