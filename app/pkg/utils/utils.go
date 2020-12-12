/*
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package utils

import (
	"errors"
	"net"
)

var (
	ErrInvalidIP     = errors.New("invalid ip")
	ErrEmptyDNSList  = errors.New("dns list is empty")
	ErrEmptyIP       = errors.New("ip is empty")
	ErrCacheNotFound = errors.New("cache is not found")
)

func ValidateIP(ip string) (string, error) {
	if _, _, err := net.SplitHostPort(ip); err != nil {
		if v := net.ParseIP(ip); v != nil {
			return net.JoinHostPort(ip, "53"), nil
		}
		return "", ErrInvalidIP
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
