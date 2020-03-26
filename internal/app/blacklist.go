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
	"net"

	"github.com/miekg/dns"
)

func (a *App) blacklist() {
	for _, ip := range a.config.BlacklistIP {
		if _, n, err := net.ParseCIDR(ip); err == nil {
			a.blacklistIPNet = append(a.blacklistIPNet, n)
		} else {
			a.blacklistIP = append(a.blacklistIP, net.ParseIP(ip))
		}
	}
}

func (a *App) InBlacklist(ip net.IP) bool {
	for _, item := range a.blacklistIP {
		if item.Equal(ip) {
			return true
		}
	}

	for _, item := range a.blacklistIPNet {
		if item.Contains(ip) {
			return true
		}
	}

	return false
}

func (a *App) BlackHole(name string) (dns.RR, error) {
	if len(a.config.BlackholeIP) == 0 {
		return nil, ErrorEmptyIP
	}
	return a.makeA(name, a.config.BlackholeIP), nil
}
