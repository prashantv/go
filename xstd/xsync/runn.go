package xsync

import (
	"iter"
	"sync"
)

// FuncN runs the given function n times in parallel.
func FuncN[T any](n int, fn func(i int) T) []T {
	var wg sync.WaitGroup
	wg.Add(n)

	res := make([]T, n)
	for i := range res {
		go func() {
			defer wg.Done()

			res[i] = fn(i)
		}()
	}
	return res
}

// RunIter runs the given function in parallel.
func RunIter[V, R any](seq iter.Seq[V], fn func(V) R) []R {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var rs []R
	set := func(i int, y R) {
		mu.Lock()
		defer mu.Unlock()

		var zero R
		for i <= len(rs) {
			rs = append(rs, zero)
		}
		rs[i] = y
	}

	var i int
	for v := range seq {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			set(i, fn(v))
		}(i)
		i++
	}

	return rs
}

// RunIter runs the given function in parallel.
func RunIter2[K, V, R any](seq iter.Seq2[K, V], fn func(K, V) R) []R {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var rs []R
	set := func(i int, y R) {
		mu.Lock()
		defer mu.Unlock()

		var zero R
		for i <= len(rs) {
			rs = append(rs, zero)
		}
		rs[i] = y
	}

	var i int
	for k, v := range seq {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			set(i, fn(k, v))
		}(i)
		i++
	}

	return rs
}
