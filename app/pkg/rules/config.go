/*
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package rules

type (
	Config struct {
		Rules []Item `yaml:"rules"`
	}
	Item struct {
		Rule string `yaml:"rule"`
		Type string `yaml:"type"`
		IP4  string `yaml:"ip4"`
		IP6  string `yaml:"ip6"`
	}
)
