package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/osspkg/go-sdk/app"
	"github.com/osspkg/go-sdk/routine"
)

type Record struct {
	values []string
	ttl    uint32
}

func (v Record) Values() []string {
	return v.values
}

func (v Record) Ttl() uint32 {
	return v.ttl
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Records struct {
	data map[string]Record
	mux  sync.RWMutex
}

func NewRecords() *Records {
	return &Records{
		data: make(map[string]Record, 1000),
	}
}

func (v *Records) Up(ctx app.Context) error {
	routine.Interval(ctx.Context(), time.Minute*15, func(ctx context.Context) {
		v.mux.Lock()
		defer v.mux.Unlock()

		currTime := uint32(time.Now().Unix())
		for name, record := range v.data {
			if record.ttl != 0 && record.ttl <= currTime {
				delete(v.data, name)
			}
		}
	})
	return nil
}

func (v *Records) Down() error {
	return nil
}

func (v *Records) Clean() {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.data = make(map[string]Record, 1000)
}

func (v *Records) ckey(rrtype uint16, name string) string {
	return fmt.Sprintf("%s:%d", name, rrtype)
}

func (v *Records) Set(rrtype uint16, name string, ttl uint32, values ...string) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.data[v.ckey(rrtype, name)] = Record{
		values: values,
		ttl:    ttl,
	}
}

//func (v *Records) Has(rrtype uint16, name string) bool {
//	v.mux.RLock()
//	defer v.mux.RUnlock()
//
//	_, ok := v.data[v.ckey(rrtype, name)]
//	return ok
//}

func (v *Records) Get(rrtype uint16, name string) *Record {
	v.mux.RLock()
	defer v.mux.RUnlock()

	if vv, ok := v.data[v.ckey(rrtype, name)]; ok {
		return &vv
	}
	return nil
}

func (v *Records) Del(rrtype uint16, name string) {
	v.mux.Lock()
	defer v.mux.Unlock()

	delete(v.data, v.ckey(rrtype, name))
}
