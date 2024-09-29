/*
 *  Copyright (c) 2020-2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ips

import (
	"net"
)

func NormalizeDNS(ips ...string) []string {
	result := make([]string, 0, len(ips))
	for _, ip := range ips {
		host, port, err := net.SplitHostPort(ip)
		if err != nil {
			host = ip
			port = "53"
		}
		if !IsValidIP(host) {
			continue
		}
		if port == "0" {
			port = "53"
		}
		result = append(result, net.JoinHostPort(host, port))
	}

	return result
}

func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}
