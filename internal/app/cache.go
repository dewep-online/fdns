/*
 * Copyright (c) 2020.  Mikhail Knyazhev <markus621@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/gpl-3.0.html>.
 */

package app

import (
	"sync"
)

type cacheMap map[string]*RuleIP

type DNSCache struct {
	config *ConfigApp
	cache  cacheMap
	sync.RWMutex

	del chan string
}

func NewDNSCache(c *ConfigApp) *DNSCache {
	dc := &DNSCache{
		config: c,
		cache:  make(cacheMap),
		del:    make(chan string, 500),
	}

	return dc
}

func (cd *DNSCache) remove() {
	for name := range cd.del {
		cd.Lock()
		delete(cd.cache, name)
		cd.Unlock()
	}
}

func (cd *DNSCache) Get(n string) *RuleIP {
	cd.RLock()
	defer cd.RUnlock()

	if rr, ok := cd.cache[n]; ok {
		return rr
	}

	return nil
}

func (cd *DNSCache) Set(n string, r *RuleIP) {
	cd.Lock()
	cd.cache[n] = r
	cd.Unlock()
}

func (cd *DNSCache) Update(name, ip4, ip6 string) {
	r := cd.Get(name)
	if r != nil && r.IsUpdatable() {
		r.SetIP4(ip4)
		r.SetIP6(ip6)
		return
	}
	nr := NewRuleIP(ip4, ip6, true)
	go nr.AutoRemove(name, cd.del)
	cd.Set(name, nr)
}
