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
	"time"
)

type RuleIP struct {
	ip4   string
	ip6   string
	canup bool
	sync.RWMutex
}

func NewRuleIP(ip4, ip6 string, canup bool) *RuleIP {
	return &RuleIP{ip4: ip4, ip6: ip6, canup: canup}
}

func (rip *RuleIP) IsUpdatable() bool {
	return rip.canup
}

func (rip *RuleIP) AutoRemove(name string, del chan string) {
	<-time.After(cTTL * time.Second)
	del <- name
}

func (rip *RuleIP) SetIP4(v string) {
	if len(v) == 0 {
		return
	}
	rip.Lock()
	rip.ip4 = v
	rip.Unlock()
}

func (rip *RuleIP) GetIP4() string {
	rip.RLock()
	defer rip.RUnlock()
	return rip.ip4
}

func (rip *RuleIP) IsEmptyIP4() bool {
	rip.RLock()
	defer rip.RUnlock()
	return len(rip.ip4) == 0
}

func (rip *RuleIP) GetIP6() string {
	rip.RLock()
	defer rip.RUnlock()
	return rip.ip6
}

func (rip *RuleIP) SetIP6(v string) {
	if len(v) == 0 {
		return
	}
	rip.Lock()
	rip.ip6 = v
	rip.Unlock()
}

func (rip *RuleIP) IsEmptyIP6() bool {
	rip.RLock()
	defer rip.RUnlock()
	return len(rip.ip6) == 0
}
