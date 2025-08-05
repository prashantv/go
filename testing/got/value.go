package got

import (
	"fmt"
	"testing"
)

type TryAssertion[V any] struct {
	v   V
	err error
}

func Try[V any](v V, err error) TryAssertion[V] {
	return TryAssertion[V]{v, err}
}

func (a TryAssertion[V]) Must() V {
	if a.err != nil {
		panic(fmt.Errorf("Must error: %w", a.err))
	}
	return a.v
}

func (ve TryAssertion[T]) Val(t testing.TB) T {
	if ve.err == nil {
		return ve.v
	}

	t.Helper()
	t.Fatalf(fmt.Sprintf("Val error: %v", ve.err))

	// Deal with "missing return" error.
	panic("Fatalf should exit goroutine")
}
