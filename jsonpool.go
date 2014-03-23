package bytepool

// Package bytepool provides a pool of fixed-length []byte

import (
	"sync/atomic"
)

type JsonPool struct {
	misses   int64
	maxBytes int64
	taken    int64
	maxTaken int64
	capacity int
	list     chan *JsonItem
	stats    map[string]int64
}

func NewJson(count int, capacity int) *JsonPool {
	p := &JsonPool{
		capacity: capacity,
		list:     make(chan *JsonItem, count),
		stats:    map[string]int64{"misses": 0, "max": 0},
	}
	for i := 0; i < count; i++ {
		p.list <- NewJsonItem(capacity, p)
	}
	return p
}

func (pool *JsonPool) Checkout() *JsonItem {
	var item *JsonItem
	select {
	case item = <-pool.list:
		taken := atomic.AddInt64(&pool.taken, 1)
		if taken > atomic.LoadInt64(&pool.maxTaken) {
			atomic.StoreInt64(&pool.maxTaken, taken)
		}
	default:
		atomic.AddInt64(&pool.misses, 1)
		item = NewJsonItem(pool.capacity, nil)
	}
	return item
}

func (pool *JsonPool) Len() int {
	return len(pool.list)
}

func (pool *JsonPool) Misses() int64 {
	return atomic.LoadInt64(&pool.misses)
}

func (pool *JsonPool) track(length int64) {
	if length > atomic.LoadInt64(&pool.maxBytes) {
		atomic.StoreInt64(&pool.maxBytes, length)
	}
	atomic.AddInt64(&pool.taken, -1)
}

func (pool *JsonPool) Stats() map[string]int64 {
	pool.stats["misses"] = atomic.SwapInt64(&pool.misses, 0)
	pool.stats["max"] = atomic.SwapInt64(&pool.maxBytes, 0)
	pool.stats["taken"] = atomic.SwapInt64(&pool.maxTaken, 0)
	return pool.stats
}
