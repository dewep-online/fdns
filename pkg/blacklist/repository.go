package blacklist

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/dewep-online/fdns/pkg/utils"
	"github.com/deweppro/go-app/application/ctx"
	"github.com/deweppro/go-errors"
	"github.com/deweppro/go-logger"

	"github.com/dewep-online/fdns/pkg/database"
)

type Repository struct {
	db             *database.Database
	blacklistIP    map[string]net.IP
	blacklistIPNet map[string]*net.IPNet
	mux            sync.RWMutex
}

func New(db *database.Database) *Repository {
	return &Repository{
		db:             db,
		blacklistIP:    make(map[string]net.IP),
		blacklistIPNet: make(map[string]*net.IPNet),
	}
}

func (v *Repository) Up(ctx ctx.Context) error {
	var timestamp int64
	utils.Interval(ctx.Context(), time.Minute*5, func(ctx context.Context) {
		err := errors.Wrap(
			v.db.GetRulesMap(ctx, database.IP, timestamp, func(m map[string]string) error {
				ipmap := make(map[string]net.IP)
				netmap := make(map[string]*net.IPNet)

				for name, ip := range m {
					if _, n, err := net.ParseCIDR(ip); err == nil {
						netmap[name] = n
					} else {
						ipmap[name] = net.ParseIP(ip)
					}
				}

				v.mux.Lock()
				v.blacklistIP, v.blacklistIPNet = ipmap, netmap
				v.mux.Unlock()

				return nil
			}),
		)
		if err != nil {
			logger.Warnf("update rules [blacklist]: %s", utils.StringError(err))
		} else {
			v.mux.RLock()
			logger.Infof("update rules [blacklist]: ip=%d net=%d", len(v.blacklistIP), len(v.blacklistIPNet))
			v.mux.RUnlock()
		}

		timestamp = time.Now().Unix()
	})

	return nil
}

func (v *Repository) Down(_ ctx.Context) error {
	return nil
}

func (v *Repository) Has(ip net.IP) bool {
	v.mux.RLock()
	defer v.mux.RUnlock()
	vv := ip.String()
	if _, ok := v.blacklistIP[vv]; ok {
		return true
	}
	for _, item := range v.blacklistIPNet {
		if item.Contains(ip) {
			return true
		}
	}
	return false
}
