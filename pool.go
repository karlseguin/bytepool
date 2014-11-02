// Package bytepool provides a pool of []byte
package bytepool

import (
	"sync/atomic"
)

type Pool struct {
	misses   int64
	size     int
	list     chan *Bytes
	stats    map[string]int64
}

func New(size, count int) *Pool {
	pool := &Pool{
		size:     size,
		list:     make(chan *Bytes, count),
		stats:    map[string]int64{"misses": 0},
	}
	for i := 0; i < count; i++ {
		pool.list <- newPooled(pool, size)
	}
	return pool
}

func (p *Pool) Checkout() *Bytes {
	select {
	case bytes := <-p.list:
		return bytes
	default:
		atomic.AddInt64(&p.misses, 1)
		return NewBytes(p.size)
	}
}

func (p *Pool) Misses() int64 {
	return atomic.LoadInt64(&p.misses)
}

func (p *Pool) Stats() map[string]int64 {
	p.stats["misses"] = atomic.SwapInt64(&p.misses, 0)
	return p.stats
}
