// Package bytepool provides a pool of []byte
package bytepool

import (
	"sync/atomic"
)

type Pool struct {
	depleted int64
	expanded int64
	size     int
	list     chan *Bytes
	stats    map[string]int64
}

// Create a new pool. The pool contains count items. Each item allocates
// an array of size bytes (but can dynamically grow)
func New(size, count int) *Pool {
	pool := &Pool{
		size:  size,
		list:  make(chan *Bytes, count),
		stats: map[string]int64{"depleted": 0, "expanded": 0},
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
// Calling this resets the counter
func (p *Pool) Depleted() int64 {
	return atomic.SwapInt64(&p.depleted, 0)
}

// Get a count of how often we had to expand an item
// beyond the initially specified size
// Calling this resets the counter
func (p *Pool) Expanded() int64 {
	return atomic.SwapInt64(&p.expanded, 0)
}

// A map containing the "expanded" and "depleted" count
// Call this resets both counters
func (p *Pool) Stats() map[string]int64 {
	p.stats["depleted"] = p.Depleted()
	p.stats["expanded"] = p.Expanded()
	return p.stats
}

func (p *Pool) onExpand() {
	atomic.AddInt64(&p.expanded, 1)
}
