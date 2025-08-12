package shardval

import (
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkLoopSerial(b *testing.B) {
	for range b.N { //nolint:revive // intentional empty block
	}
}

func BenchmarkCounterLowBound(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			i++
		}
		_ = i
	})
}

func BenchmarkCounterAtomic(b *testing.B) {
	var i atomic.Int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i.Add(1)
		}
	})
}

func BenchmarkCounterMutex(b *testing.B) {
	var mu sync.Mutex
	var i int64
	for b.Loop() {
		mu.Lock()
		i++
		mu.Unlock()
	}
}

func BenchmarkCounterSharded(b *testing.B) {
	// Since there's no synchronization possible with `*int64`,
	// `ForEach` cannot be called till the `With` calls are complete.
	// This is useful when the results are not necessary till the work is complete.
	var sharded Value[int64]
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sharded.With(func(v *int64) {
				*v++
			})
		}
	})
}

func BenchmarkCounterShardedAtomic(b *testing.B) {
	// Atomics allow concurrent use, so `ForAll` can be called concurrently to `With`.
	var sharded Value[atomic.Int64]
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sharded.With(func(v *atomic.Int64) {
				v.Add(1)
			})
		}
	})
}

func BenchmarkCounterShardedMutex(b *testing.B) {
	// A mutex should be used in scenarios where `With` and `ForAll` need to be
	// called concurrently, and atomics (or similar) aren't available.
	type lockedInt struct {
		mu sync.Mutex
		v  int
	}
	var sharded Value[lockedInt]
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sharded.With(func(v *lockedInt) {
				v.mu.Lock()
				defer v.mu.Unlock()
				v.v++
			})
		}
	})
}
