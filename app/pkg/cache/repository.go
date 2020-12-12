/*
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package cache

import (
	"sync"
	"time"

	"fdns/app/pkg/utils"
)

type (
	Cache      map[string]*Item
	Repository struct {
		data Cache
		sync.RWMutex
	}
)

func New() *Repository {
	return &Repository{
		data: make(Cache),
	}
}

func (o *Repository) Get(name string) *Item {
	o.RLock()
	defer o.RUnlock()
	if d, ok := o.data[name]; ok {
		return d
	}
	return nil
}

func (o *Repository) Set(name string, d *Item) {
	o.Lock()
	o.data[name] = d
	o.Unlock()
}

func (o *Repository) Del(name string) {
	o.Lock()
	delete(o.data, name)
	o.Unlock()
}

func (o *Repository) Update(name string, ip4, ip6 []string) {
	d := o.Get(name)
	if d != nil {
		d.SetIP4(ip4...)
		d.SetIP6(ip6...)
	} else {
		nr := Create(ip4, ip6, true)
		o.Set(name, nr)
		time.AfterFunc(utils.TTL*time.Second, func() {
			o.Del(name)
		})
	}
}
