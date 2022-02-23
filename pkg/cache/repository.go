package cache

import (
	"sync"
	"time"

	"github.com/deweppro/go-app/application"
)

type Repository struct {
	data map[string]*Record
	temp map[string]*Record

	ctx *application.ForceClose
	wg  sync.WaitGroup
	mux sync.RWMutex
}

func New(ctx *application.ForceClose) *Repository {
	return &Repository{
		data: make(map[string]*Record),
		temp: make(map[string]*Record),
		ctx:  ctx,
	}
}

func (v *Repository) Get(name string) *Record {
	v.mux.RLock()
	defer v.mux.RUnlock()

	if d, ok := v.data[name]; ok {
		return d
	}
	if d, ok := v.temp[name]; ok {
		return d
	}
	return nil
}

func (v *Repository) Set(name string, ip4, ip6 []string, ttl int64) {
	v.mux.Lock()
	defer v.mux.Unlock()

	m := NewRecord(ip4, ip6, ttl)
	if m.IsStatic() {
		v.temp[name] = m
	} else {
		v.data[name] = m
	}
}

func (v *Repository) Del(name string) {
	v.mux.Lock()
	defer v.mux.Unlock()

	delete(v.data, name)
}

func (v *Repository) List(call func(name string, ip []string)) {
	v.mux.RLock()
	defer v.mux.RUnlock()

	for name, vv := range v.data {
		ip := vv.AllIPs()
		call(name, ip)
	}
}

func (v *Repository) Up() error {
	v.wg.Add(1)
	go func() {
		defer v.wg.Done()

		timer := time.NewTicker(time.Minute)
		defer timer.Stop()

		for {
			select {
			case <-v.ctx.C.Done():
				return

			case t := <-timer.C:
				v.mux.Lock()
				for name, record := range v.data {
					if record.GetTTL() <= t.Unix() {
						delete(v.data, name)
					}
				}
				v.mux.Unlock()
			}
		}
	}()
	return nil
}

func (v *Repository) Down() error {
	v.wg.Wait()
	return nil
}
