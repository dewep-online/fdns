package api

import (
	"net/http"

	"github.com/dewep-games/fdns/pkg/cache"

	"github.com/deweppro/go-http/web/routes"
)

//API model
type API struct {
	route *routes.Router
	repo  *cache.Repository
}

//NewAPI init api
func NewAPI(route *routes.Router, repo *cache.Repository) *API {
	return &API{
		route: route,
		repo:  repo,
	}
}

//Up startup api service
func (v *API) Up() error {
	v.route.Route("/api", v.Index, http.MethodGet)
	v.route.Route("/api/cache/list", v.CacheList, http.MethodGet)

	return nil
}

//Down shutdown api service
func (v *API) Down() error {
	return nil
}
