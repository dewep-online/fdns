package webserver

import (
	"github.com/dewep-online/fdns/pkg/cache"
	"github.com/dewep-online/fdns/pkg/database"
	"github.com/dewep-online/fdns/pkg/dnscli"
	"github.com/dewep-online/fdns/pkg/rules"
	"github.com/deweppro/go-http/pkg/routes"
	"github.com/deweppro/go-logger"
	"github.com/deweppro/go-static"
)

//WebServer model
type WebServer struct {
	conf  *MiddlewareConfig
	route *routes.Router
	cache *static.Cache
	repo  *cache.Repository
	rules *rules.Repository
	cli   *dnscli.Client
	db    *database.Database
}

func New(route *routes.Router, conf *MiddlewareConfig, repo *cache.Repository,
	rules *rules.Repository, db *database.Database, cli *dnscli.Client) *WebServer {
	return &WebServer{
		cache: static.New(),
		route: route,
		conf:  conf,
		repo:  repo,
		rules: rules,
		cli:   cli,
		db:    db,
	}
}

//Up startup api service
func (v *WebServer) Up() error {
	v.route.Global(routes.RecoveryMiddleware(logger.Default()))
	v.route.Global(routes.ThrottlingMiddleware(v.conf.Middleware.Throttling))

	v.RegisterUI()
	v.RegisterAPI()

	return nil
}

//Down shutdown api service
func (v *WebServer) Down() error {
	return nil
}
