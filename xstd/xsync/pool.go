package xsync

import "sync"

// Pool is a generic wrapper for [sync.Pool].
type Pool[T any] struct {
	pool sync.Pool

	New func() T
}

// Get returns an element from the underlying pool, see [sync.Pool.Get].
// If the pool returns no element, p.New is returned if set, or a zero value is returned.
func (p *Pool[T]) Get() T {
	x := p.pool.Get()
	if x == nil {
		if p.New == nil {
			var zero T
			return zero
		}
		return p.New()
	}

	return x.(T)
}

// Put adds x to the pool.
func (p *Pool[T]) Put(x T) {
	p.pool.Put(x)
}
