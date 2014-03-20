package bytepool

// Package bytepool provides a pool of fixed-length []byte

import (
	"sync/atomic"
)

type Pool struct {
	misses   int64
	maxBytes int64
	capacity int
	list     chan *Item
	stats    map[string]int64
}

func New(count int, capacity int) *Pool {
	pool := &Pool{
		capacity: capacity,
		list:     make(chan *Item, count),
		stats:    map[string]int64{"misses": 0, "max": 0},
	}
	pool.populate()
	return pool
}

// not thread safe, just here to make fluent-configs a little cleaner
func (pool *Pool) SetCapacity(capacity int) {
	pool.capacity = capacity
	pool.list = make(chan *Item, len(pool.list))
	pool.populate()
}

// not thread safe, just here to make fluent-configs a little cleaner
func (pool *Pool) SetCount(count int) {
	pool.list = make(chan *Item, count)
	pool.populate()
}

func (pool *Pool) Checkout() *Item {
	var item *Item
	select {
	case item = <-pool.list:
	default:
		atomic.AddInt64(&pool.misses, 1)
		item = newItem(pool.capacity, nil)
	}
	return item
}

func (pool *Pool) Len() int {
	return len(pool.list)
}

func (pool *Pool) Misses() int64 {
	return atomic.LoadInt64(&pool.misses)
}

func (pool *Pool) populate() {
	for i := 0; i < cap(pool.list); i++ {
		pool.list <- newItem(pool.capacity, pool)
	}
}

func (pool *Pool) track(length int64) {
	if length > atomic.LoadInt64(&pool.maxBytes) {
		atomic.StoreInt64(&pool.maxBytes, length)
	}
}

func (pool *Pool) Stats() map[string]int64 {
	pool.stats["misses"] = atomic.SwapInt64(&pool.misses, 0)
	pool.stats["max"] = atomic.SwapInt64(&pool.maxBytes, 0)
	return pool.stats
}
