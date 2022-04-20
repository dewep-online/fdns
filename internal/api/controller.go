package api

import (
	"net/http"
	"sort"

	"github.com/deweppro/go-http/pkg/httputil/enc"
)

//go:generate easyjson

//easyjson:json
type (
	CacheList     []CacheListItem
	CacheListItem struct {
		Domain string   `json:"domain"`
		IP     []string `json:"ip"`
	}
)

//Index controller
func (v *API) CacheList(w http.ResponseWriter, r *http.Request) {
	list := make(CacheList, 0)
	v.repo.List(func(name string, ip []string) {
		list = append(list, CacheListItem{
			Domain: name,
			IP:     ip,
		})
	})
	sort.Slice(list, func(i, j int) bool {
		return list[i].Domain < list[j].Domain
	})
	enc.JSON(w, list)
}
