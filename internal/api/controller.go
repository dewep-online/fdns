package api

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/deweppro/go-logger"
)

//go:generate easyjson

//Index controller
func (v *API) Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Hello")); err != nil {
		logger.Errorf("Index: %s", err.Error())
	}
}

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
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		logger.Errorf("cache-list: %s", err.Error())
	}
}
