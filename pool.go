// Package bytepool provides a pool of []byte
package bytepool

import (
	"sync/atomic"
)

type Pool struct {
	depleted   int64
	size     int
	list     chan *Bytes
	stats    map[string]int64
}

// Create a new pool. The pool contains count items. Each item allocates
// an array of size bytes (but can dynamically grow)
func New(size, count int) *Pool {
	pool := &Pool{
		size:     size,
		list:     make(chan *Bytes, count),
		stats:    map[string]int64{"depleted": 0},
	}
	for i := 0; i < count; i++ {
		pool.list <- newPooled(pool, size)
	}
	return pool
}

// Get an item from the pool
func (p *Pool) Checkout() *Bytes {
	select {
	case bytes := <-p.list:
		return bytes
	default:
		atomic.AddInt64(&p.depleted, 1)
		return NewBytes(p.size)
	}
}

// Get a count of how often Checkout() was called
// but no item was available (thus causing an item to be
// created on the fly)
func (p *Pool) Depleted() int64 {
	return atomic.SwapInt64(&p.depleted, 0)
}

// Same as Depleted, but returned as a map
func (p *Pool) Stats() map[string]int64 {
	p.stats["depleted"] = p.Depleted()
	return p.stats
}
