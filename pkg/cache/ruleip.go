package cache

import (
	"sync"
)

type (
	Item struct {
		ip4 itemMap
		ip6 itemMap
		is  bool
		sync.RWMutex
	}
	itemMap map[string]struct{}
)

func Create(ip4, ip6 []string, canup bool) *Item {
	item := &Item{ip4: make(itemMap), ip6: make(itemMap), is: true}
	item.SetIP4(ip4...)
	item.SetIP6(ip6...)
	item.SetUpdatable(canup)
	return item
}

func (o *Item) IsUpdatable() bool {
	o.RLock()
	defer o.RUnlock()
	return o.is
}

func (o *Item) SetUpdatable(v bool) {
	o.Lock()
	o.is = v
	o.Unlock()
}

func (o *Item) SetIP4(v ...string) {
	if len(v) == 0 || !o.IsUpdatable() {
		return
	}
	o.Lock()
	for _, ip := range v {
		o.ip4[ip] = struct{}{}
	}
	o.Unlock()
}

func (o *Item) GetIP4() (v []string) {
	o.RLock()
	defer o.RUnlock()
	for ip := range o.ip4 {
		v = append(v, ip)
	}
	return
}

func (o *Item) HasIP4() bool {
	o.RLock()
	defer o.RUnlock()
	return len(o.ip4) > 0
}

func (o *Item) GetIP6() (v []string) {
	o.RLock()
	defer o.RUnlock()
	for ip := range o.ip6 {
		v = append(v, ip)
	}
	return
}

func (o *Item) SetIP6(v ...string) {
	if len(v) == 0 || !o.IsUpdatable() {
		return
	}
	o.Lock()
	for _, ip := range v {
		o.ip6[ip] = struct{}{}
	}
	o.Unlock()
}

func (o *Item) HasIP6() bool {
	o.RLock()
	defer o.RUnlock()
	return len(o.ip6) > 0
}
