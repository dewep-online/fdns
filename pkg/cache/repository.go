package cache

import (
	"strings"
	"sync"
	"time"

	"github.com/deweppro/go-app/application/ctx"
)

type Repository struct {
	dyn map[string]*Record
	fix map[string]*Record

	wg  sync.WaitGroup
	mux sync.RWMutex
}

func New() *Repository {
	return &Repository{
		dyn: make(map[string]*Record),
		fix: make(map[string]*Record),
	}
}

func (v *Repository) Get(name string) *Record {
	v.mux.RLock()
	defer v.mux.RUnlock()

	if d, ok := v.fix[name]; ok {
		return d
	}
	if d, ok := v.dyn[name]; ok {
		return d
	}
	return nil
}

func (v *Repository) Set(name string, ip4, ip6 []string, ttl int64) {
	v.mux.Lock()
	defer v.mux.Unlock()

	m := NewRecord(ip4, ip6, ttl)
	if m.IsStatic() {
		v.fix[name] = m
	} else {
		v.dyn[name] = m
	}
}

func (v *Repository) DelDynamic(name string) {
	v.mux.Lock()
	defer v.mux.Unlock()

	delete(v.dyn, name)
}

func (v *Repository) DelFixed(name string) {
	v.mux.Lock()
	defer v.mux.Unlock()

	delete(v.fix, name)
}

func (v *Repository) DelByCallback(call func(name string) bool) {
	var result []string
	v.mux.RLock()
	for name, _ := range v.fix {
		if call(name) {
			result = append(result, name)
		}
	}
	for name, _ := range v.dyn {
		if call(name) {
			result = append(result, name)
		}
	}
	v.mux.RUnlock()

	v.mux.Lock()
	for _, name := range result {
		delete(v.fix, name)
		delete(v.dyn, name)
	}
	v.mux.Unlock()
}

func (v *Repository) Reset() {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.dyn = make(map[string]*Record)
	v.fix = make(map[string]*Record)
}

func (v *Repository) List(dyn bool, filter string, call func(name string, ip []string, ttl string)) {
	v.mux.RLock()
	defer v.mux.RUnlock()

	now := time.Now()

	fn := func(list map[string]*Record) {
		for name, vv := range list {
			if len(filter) > 2 && strings.Index(name, filter) == -1 {
				continue
			}
			ip := vv.AllIPs()
			var ttl string
			if vv.ttl > 0 {
				ttl = time.Unix(vv.ttl, 0).Sub(now).String()
			}
			call(strings.Trim(name, "."), ip, ttl)
		}
	}

	if dyn {
		fn(v.dyn)
	} else {
		fn(v.fix)
	}
}

func (v *Repository) Up(ctx ctx.Context) error {
	v.wg.Add(1)
	go func() {
		defer v.wg.Done()

		timer := time.NewTicker(time.Minute)
		defer timer.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case t := <-timer.C:
				v.mux.Lock()
				for name, record := range v.dyn {
					if record.GetTTL() <= t.Unix() {
						delete(v.dyn, name)
					}
				}
				v.mux.Unlock()
			}
		}
	}()
	return nil
}

func (v *Repository) Down(_ ctx.Context) error {
	v.wg.Wait()
	return nil
}
