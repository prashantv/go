package xslices

import (
	"context"
	"fmt"
)

// Map runs the mapping function to map a slice to a new slice.
func Map[X, Y any](xs []X, fn func(X) Y) []Y {
	if xs == nil {
		return nil
	}

	ys := make([]Y, len(xs))
	// Iterate by index only, so xs[i] is directly read into the function argument.
	for i := range xs {
		ys[i] = fn(xs[i])
	}
	return ys
}

// MapErr runs the mapping function to map a slice to a new slice.
// The mapping function may return an error which stops and returns the partially mapped slice.
func MapErr[X, Y any](xs []X, fn func(X) (Y, error)) ([]Y, error) {
	if xs == nil {
		return nil, nil
	}

	ys := make([]Y, len(xs))
	for i := range xs {
		var err error
		ys[i], err = fn(xs[i])
		if err != nil {
			return nil, fmt.Errorf("index %d: %w", i, err)
		}
	}

	return ys, nil
}

// MapCtx runs the context-aware mapping function to map a slice to a new slice.
// The mapping function may return an error which stops and returns the partially mapped slice.
// MapCtx does not explicitly check for context errors, the mapping function is expected to respect and propagate context errors.
func MapCtx[X, Y any](ctx context.Context, xs []X, fn func(context.Context, X) (Y, error)) ([]Y, error) {
	if xs == nil {
		return nil, nil
	}

	ys := make([]Y, len(xs))
	for i := range xs {
		var err error
		ys[i], err = fn(ctx, xs[i])
		if err != nil {
			return nil, fmt.Errorf("index %d: %w", i, err)
		}
	}
	return ys, nil
}
