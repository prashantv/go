package xsync

// RunN runs the given function in parallel and returns
func RunN[T any](n int, fn func(i int) T) []T {
	res := make([]T, n)
	for i := range res {
		res[i] = fn(i)
	}
	return res
}
