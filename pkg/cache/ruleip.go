package cache

import "strings"

type Record struct {
	ip4 map[string]struct{}
	ip6 map[string]struct{}
	ttl int64
}

func NewRecord(ip4, ip6 []string, ttl int64) *Record {
	item := &Record{
		ip4: make(map[string]struct{}),
		ip6: make(map[string]struct{}),
		ttl: ttl,
	}
	for _, ip := range ip4 {
		item.ip4[ip] = struct{}{}
	}
	for _, ip := range ip6 {
		item.ip6[ip] = struct{}{}
	}
	return item
}

func (v *Record) IsStatic() bool {
	return v.ttl == 0
}

func (v *Record) GetTTL() int64 {
	return v.ttl
}

func (v *Record) GetIP4() []string {
	vv := make([]string, 0, len(v.ip4))
	for ip := range v.ip4 {
		vv = append(vv, ip)
	}
	return vv
}

func (v *Record) HasIP4() bool {
	return len(v.ip4) > 0
}

func (v *Record) GetIP6() []string {
	vv := make([]string, 0, len(v.ip6))
	for ip := range v.ip6 {
		vv = append(vv, ip)
	}
	return vv
}

func (v *Record) HasIP6() bool {
	return len(v.ip6) > 0
}

func (v *Record) AllIPs() []string {
	vv := make([]string, 0, len(v.ip4)+len(v.ip6))
	for ip := range v.ip4 {
		vv = append(vv, ip)
	}
	for ip := range v.ip6 {
		vv = append(vv, ip)
	}
	return vv
}

func (v *Record) AllIPsString() string {
	vv := make([]string, 0, len(v.ip4)+len(v.ip6))
	for ip := range v.ip4 {
		vv = append(vv, ip)
	}
	for ip := range v.ip6 {
		vv = append(vv, ip)
	}
	return strings.Join(vv, ", ")
}
