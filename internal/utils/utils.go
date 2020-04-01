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

package utils

import (
	"errors"
	"net"
)

var (
	ErrorInvalidIP     = errors.New("invalid ip")
	ErrorEmptyDNSList  = errors.New("dns list is empty")
	ErrorEmptyIP       = errors.New("ip is empty")
	ErrorCacheNotFound = errors.New("cache is not found")
)

func ValidateIP(ip string) (string, error) {
	if _, _, err := net.SplitHostPort(ip); err != nil {
		if v := net.ParseIP(ip); v != nil {
			return net.JoinHostPort(ip, "53"), nil
		}
		return "", ErrorInvalidIP
	}
	return ip, nil
}

func ValidateIPs(list []string) (result []string) {
	for _, ip := range list {
		if v, er := ValidateIP(ip); er == nil {
			result = append(result, v)
		}
	}
	return
}
