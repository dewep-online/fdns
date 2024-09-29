/*
 *  Copyright (c) 2020-2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cache

import (
	"context"
	"fmt"
	"time"

	"go.osspkg.com/ioutils/cache"
)

type Record struct {
	Value []string
	TTL   uint32
}

type Records struct {
	data cache.TCacheTTL[string, *Record]
}

func NewRecords(ctx context.Context) *Records {
	return &Records{
		data: cache.NewWithTTL[string, *Record](ctx, 15*time.Minute),
	}
}

func (v *Records) key(qtype uint16, name string) string {
	return fmt.Sprintf("%d %s", qtype, name)
}

func (v *Records) Set(qtype uint16, name string, ttl uint32, values ...string) {
	v.data.SetWithTTL(
		v.key(qtype, name),
		&Record{
			Value: nil,
			TTL:   ttl,
		},
		time.Unix(int64(ttl), 0),
	)
}

func (v *Records) Has(qtype uint16, name string) bool {
	return v.data.Has(v.key(qtype, name))
}

func (v *Records) Get(qtype uint16, name string) (*Record, bool) {
	return v.data.Get(v.key(qtype, name))
}

func (v *Records) Del(qtype uint16, name string) {
	v.data.Del(v.key(qtype, name))
}
