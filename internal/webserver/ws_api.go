package webserver

//go:generate easyjson

import (
	"net/http"
	"sort"
	"strings"

	"github.com/deweppro/go-http/pkg/httputil/dec"
	"github.com/deweppro/go-http/pkg/httputil/enc"

	"github.com/dewep-online/fdns/pkg/database"

	"github.com/dewep-online/fdns/pkg/utils"

	"github.com/deweppro/go-logger"
)

func (v *WebServer) RegisterAPI() {
	v.route.Route("/api", v.Index, http.MethodGet)
	v.route.Route("/api/cache/list", v.CacheList, http.MethodGet)
	v.route.Route("/api/cache/block", v.BlockDomain, http.MethodPost)

	v.route.Route("/api/adblock/list/uri", v.AdblockURIList, http.MethodGet)
	v.route.Route("/api/adblock/list/domain", v.AdblockDomainList, http.MethodGet)
	v.route.Route("/api/adblock/active", v.AdblockActive, http.MethodPost)

	v.route.Route("/api/fixed/list", v.FixedList, http.MethodGet)
	v.route.Route("/api/fixed/save", v.FixedSave, http.MethodPost)
	v.route.Route("/api/fixed/active", v.FixedActive, http.MethodPost)
	v.route.Route("/api/fixed/delete", v.FixedDelete, http.MethodPost)
}

func (v *WebServer) Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Hello")); err != nil {
		logger.Errorf("Index: %s", err.Error())
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//easyjson:json
type (
	CacheList     []CacheListItem
	CacheListItem struct {
		Domain string   `json:"domain"`
		IP     []string `json:"ip"`
		TTL    string   `json:"ttl"`
	}
)

func (v *WebServer) CacheList(w http.ResponseWriter, r *http.Request) {
	var dyn bool
	t := r.URL.Query().Get("type")
	filter := r.URL.Query().Get("filter")
	if t == "1" {
		dyn = true
	}

	list := make(CacheList, 0)
	v.repo.List(dyn, strings.TrimSpace(filter), func(name string, ip []string, ttl string) {
		list = append(list, CacheListItem{
			Domain: name,
			IP:     ip,
			TTL:    ttl,
		})
	})
	sort.Slice(list, func(i, j int) bool {
		return list[i].Domain < list[j].Domain
	})
	enc.JSON(w, list)
}

//easyjson:json
type BlockDomainModel struct {
	Domain string `json:"domain"`
}

func (v *WebServer) BlockDomain(w http.ResponseWriter, r *http.Request) {
	var err error
	mod := BlockDomainModel{}
	if err = dec.JSON(r, &mod); err != nil {
		enc.Error(w, err)
		return
	}
	mod.Domain, err = utils.ValidateDomain(mod.Domain)
	if err != nil {
		enc.Error(w, err)
		return
	}

	err = v.db.SetRules(r.Context(), database.Host, mod.Domain, "", database.ActiveTrue)
	if err != nil {
		enc.Error(w, err)
		return
	}

	v.repo.Set(mod.Domain, nil, nil, 0)
	v.repo.DelDynamic(mod.Domain)
	w.WriteHeader(http.StatusOK)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//easyjson:json
type (
	AdblockDomainModel struct {
		Tag    string `json:"tag"`
		Domain string `json:"domain"`
		Active bool   `json:"active"`
	}
	AdblockURIModel struct {
		Tag    string `json:"tag"`
		URI    string `json:"uri"`
		Active bool   `json:"active"`
	}
)

func (v *WebServer) AdblockURIList(w http.ResponseWriter, r *http.Request) {
	vv, err := v.db.GetBlacklistURI(r.Context(), 0)
	if err != nil {
		enc.Error(w, err)
		return
	}
	list := make([]AdblockURIModel, 0, len(vv))
	for _, model := range vv {
		list = append(list, AdblockURIModel{
			Tag:    model.Tag,
			URI:    model.URI,
			Active: model.Active == database.ActiveTrue,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].URI < list[j].URI
	})
	enc.JSON(w, list)
}

func (v *WebServer) AdblockDomainList(w http.ResponseWriter, r *http.Request) {
	vv, err := v.db.GetBlacklistDomain(r.Context(), 0)
	if err != nil {
		enc.Error(w, err)
		return
	}
	list := make([]AdblockDomainModel, 0, len(vv))
	for _, model := range vv {
		list = append(list, AdblockDomainModel{
			Tag:    model.Tag,
			Domain: model.Domain,
			Active: model.Active == database.ActiveTrue,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Domain < list[j].Domain
	})
	enc.JSON(w, list)
}

//easyjson:json
type (
	AdblockActiveModel struct {
		Domain string `json:"domain"`
		Active bool   `json:"active"`
	}
)

func (v *WebServer) AdblockActive(w http.ResponseWriter, r *http.Request) {
	var err error
	mod := AdblockActiveModel{}
	if err = dec.JSON(r, &mod); err != nil {
		enc.Error(w, err)
		return
	}
	mod.Domain, err = utils.ValidateDomain(mod.Domain)
	if err != nil {
		enc.Error(w, err)
		return
	}

	active := database.ActiveFalse
	if mod.Active {
		active = database.ActiveTrue
	}

	if err = v.db.SetBlacklistDomainActive(r.Context(), mod.Domain, active); err != nil {
		enc.Error(w, err)
		return
	}

	if mod.Active {
		v.repo.DelDynamic(mod.Domain)
		v.repo.Set(mod.Domain, nil, nil, 0)
	} else {
		v.repo.DelFixed(mod.Domain)
	}

	w.WriteHeader(http.StatusOK)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//easyjson:json
type (
	FixedList     []FixedListItem
	FixedListItem struct {
		Types  string `json:"types"`
		Origin string `json:"origin"`
		Domain string `json:"domain"`
		IPs    string `json:"ips"`
		Active bool   `json:"active"`
	}
)

func (v *WebServer) FixedList(w http.ResponseWriter, r *http.Request) {
	rules, err := v.db.GetAllRules(r.Context(), 0)
	if err != nil {
		enc.Error(w, err)
		return
	}
	list := make(FixedList, 0, len(rules))
	for _, rule := range rules {
		list = append(list, FixedListItem{
			Types:  rule.Types,
			Origin: rule.Domain,
			Domain: rule.Domain,
			IPs:    rule.IPs,
			Active: rule.Active == database.ActiveTrue,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Types < list[j].Types
	})
	enc.JSON(w, list)
}

func (v *WebServer) FixedSave(w http.ResponseWriter, r *http.Request) {
	var err error
	mod := FixedListItem{}
	if err = dec.JSON(r, &mod); err != nil {
		enc.Error(w, err)
		return
	}

	tt, err := database.ValidateType(mod.Types)
	if err != nil {
		enc.Error(w, err)
		return
	}

	if tt == database.Host {
		mod.Domain, err = utils.ValidateDomain(mod.Domain)
		if err != nil {
			enc.Error(w, err)
			return
		}
	}

	active := database.ActiveFalse
	if mod.Active {
		active = database.ActiveTrue
	}

	if err = v.db.SetRules(r.Context(), tt, mod.Domain, mod.IPs, active); err != nil {
		enc.Error(w, err)
		return
	}

	switch tt {
	case database.Host:
		if mod.Active {
			ip4, ip6 := utils.DecodeIPs(mod.IPs)
			v.repo.DelDynamic(mod.Domain)
			v.repo.Set(mod.Domain, ip4, ip6, 0)
		} else {
			v.repo.DelFixed(mod.Domain)
		}

	case database.DNS, database.Regex, database.Query:
		if mod.Active {
			v.rules.ReplaceRexResolve(tt, mod.Origin, mod.Domain, mod.IPs)
		} else {
			v.rules.DeleteRexResolve(mod.Origin)
		}

	case database.NS:
		v.cli.UpgradeNS(r.Context())
	}

	mod.Origin = mod.Domain

	enc.JSON(w, &mod)
}

func (v *WebServer) FixedActive(w http.ResponseWriter, r *http.Request) {
	v.FixedSave(w, r)
}

func (v *WebServer) FixedDelete(w http.ResponseWriter, r *http.Request) {
	var err error
	mod := FixedListItem{}
	if err = dec.JSON(r, &mod); err != nil {
		enc.Error(w, err)
		return
	}

	tt, err := database.ValidateType(mod.Types)
	if err != nil {
		enc.Error(w, err)
		return
	}

	if err = v.db.DelRules(r.Context(), tt, mod.Origin); err != nil {
		enc.Error(w, err)
		return
	}

	switch tt {
	case database.Host:
		v.repo.DelFixed(mod.Origin)
	case database.DNS, database.Regex, database.Query:
		v.rules.DeleteRexResolve(mod.Origin)
	case database.NS:
		v.cli.UpgradeNS(r.Context())
	}

	w.WriteHeader(http.StatusOK)
}
