package bytepool

// Package bytepool provides a pool of fixed-length []byte

import (
	"sync/atomic"
)

type Pool struct {
	misses   int32
	capacity int
	list     chan *Item
}

func New(count int, capacity int) *Pool {
	pool := &Pool{
		capacity: capacity,
		list:     make(chan *Item, count),
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
		atomic.AddInt32(&pool.misses, 1)
		item = newItem(pool.capacity, nil)
	}
	return item
}

func (pool *Pool) Len() int {
	return len(pool.list)
}

func (pool *Pool) Misses() int32 {
	return atomic.LoadInt32(&pool.misses)
}

func (pool *Pool) populate() {
	for i := 0; i < cap(pool.list); i++ {
		pool.list <- newItem(pool.capacity, pool)
	}
}
