package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/alecthomas/binary"
	"github.com/xmdhs/creditget/db"
	"github.com/xmdhs/creditget/model"
	"golang.org/x/sync/singleflight"
)

var _ db.DB = &MemCache{}

type MemCache struct {
	c *fastcache.Cache
	db.DB
	s singleflight.Group
}

type ttlCache[V any] struct {
	TimeOut int64
	V       V
}

func NewMemCache(maxBytes int, db db.DB) *MemCache {
	c := fastcache.New(maxBytes)
	return &MemCache{c: c, DB: db}
}

func cacheGetWarp[K any](m *MemCache, get func() (K, error), key string, timeout time.Duration) (K, error) {
	r := m.c.GetBig(nil, []byte(key))
	if r != nil {
		me := ttlCache[K]{}
		err := binary.Unmarshal(r, &me)
		if err == nil && time.Unix(me.TimeOut, 0).After(time.Now()) {
			return me.V, nil
		}
	}
	v, err, _ := m.s.Do(key, func() (interface{}, error) {
		mm, err := get()
		if err != nil {
			return mm, err
		}
		return mm, nil
	})
	if err != nil {
		var mm K
		return mm, err
	}
	mm := v.(K)

	b, err := binary.Marshal(ttlCache[K]{V: mm, TimeOut: time.Now().Add(timeout).Unix()})
	if err != nil {
		panic(err)
	}
	m.c.SetBig([]byte(key), b)
	return mm, nil
}

func (m *MemCache) GetCreditInfo(cxt context.Context, uid int) (*model.CreditInfo, error) {
	return cacheGetWarp(m, func() (*model.CreditInfo, error) {
		return m.DB.GetCreditInfo(cxt, uid)
	}, strconv.Itoa(uid), 48*time.Hour)
}

func (m *MemCache) GetRank(cxt context.Context, uid int, field string) (int, error) {
	return cacheGetWarp(m, func() (int, error) {
		return m.DB.GetRank(cxt, uid, field)
	}, strconv.Itoa(uid)+field, 48*time.Hour)
}
