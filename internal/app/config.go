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

type ConfigApp struct {
	Cache       Cache    `yaml:"cache"`
	DNS         []string `yaml:"dns"`
	BlackholeIP string   `yaml:"blackholeip"`
	BlacklistIP []string `yaml:"blacklistip"`
	Rules       []Rule   `yaml:"rules"`
}

type Cache struct {
	CacheFile string `yaml:"file"`
	CacheTtl  int    `yaml:"ttl"`
}

type Rule struct {
	Rule string `yaml:"rule"`
	Type string `yaml:"type"`
	IP4  string `yaml:"ip4"`
	IP6  string `yaml:"ip6"`
}
