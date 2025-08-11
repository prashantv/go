package shardval

import (
	"runtime"
	"sync"
)

// Value is a sharded value implemented using a [sync.Pool] for per-P locality.
// It is used to represent a value that's sharded into many pieces
// which are updated in `With`, and then assembled in `ForEach`.
//
// The number of shards should scale with O(GOMAXPROCS), though there is no hard limit.
type Value[T any] struct {
	pool sync.Pool

	mu    sync.Mutex // below fields are protected by mu.
	all   []*T
	reuse []*shard[T] // GC'd elements, reused if the pool is empty.
}

// shard is an indirection used to detects elements dropped from the pool (via finalizers).
type shard[T any] struct {
	val *T
}

// With runs the given function with exclusive access to the shard.
//
// Many callers may call With concurrently.
//
// If called concurrently with [ForEach], the same value may be processed by both
// functions so addititional synchronization within T is required.
func (val *Value[T]) With(fn func(*T)) {
	s := val.get()
	defer val.pool.Put(s)

	fn(s.val)
}

// ForEach iterates over all shards.
//
// If called concurretnly with [With], the same value may be processed by both functions
// so addititional synchronization within T is required.
func (val *Value[T]) ForEach(fn func(*T)) {
	val.mu.Lock()
	defer val.mu.Unlock()

	for i := range val.all {
		fn(val.all[i])
	}
}

// Shards returns the number of shards.
func (val *Value[T]) Shards() int {
	val.mu.Lock()
	defer val.mu.Unlock()

	return len(val.all)
}

func (val *Value[T]) get() *shard[T] {
	s, ok := val.pool.Get().(*shard[T])
	if ok {
		return s
	}

	val.mu.Lock()
	defer val.mu.Unlock()

	// The lock may take time, so check if there's a free shard again.
	s, ok = val.pool.Get().(*shard[T])
	if ok {
		return s
	}

	if s, ok := val.reuseLocked(); ok {
		return s
	}

	s = val.newLocked()
	return s
}

func (val *Value[T]) reuseLocked() (*shard[T], bool) {
	if len(val.reuse) == 0 {
		return nil, false
	}

	s := val.reuse[0]
	val.reuse[0] = nil
	val.reuse = val.reuse[1:]
	return s, true
}

func (val *Value[T]) newLocked() *shard[T] {
	s := &shard[T]{
		val: new(T),
	}
	val.all = append(val.all, s.val)
	runtime.SetFinalizer(s, val.finalize)
	return s
}

func (val *Value[T]) finalize(s *shard[T]) {
	val.mu.Lock()
	defer val.mu.Unlock()

	val.reuse = append(val.reuse, s)
	runtime.SetFinalizer(s, val.finalize)
}
