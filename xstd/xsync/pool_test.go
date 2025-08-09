package xsync

import (
	"reflect"
	"testing"
)

func TestPool(t *testing.T) {
	t.Run("no New", func(t *testing.T) {
		var p Pool[*string]
		for range 100 {
			assertNil(t, p.Get())
			assertNil(t, p.Get())
		}

		for range 100 {
			s := "hello"
			p.Put(&s)
			got := p.Get()
			if got != nil {
				return
			}
			// Pool has P-local caches, so misses are exected, retry.
		}
		t.Fatal("pool.Get only returned nil")
	})

	t.Run("New", func(t *testing.T) {
		p := Pool[*string]{
			New: func() *string {
				s := "initial"
				return &s
			},
		}

		for range 100 {
			got := p.Get()
			assertNotNil(t, got)
			assertEq(t, "initial", *got)
		}

		for range 100 {
			s := "put"
			p.Put(&s)
			got := p.Get()
			assertNotNil(t, got)

			if *got == "put" {
				return
			}

			// Pool has P-local caches, so misses are exected, retry.
			assertEq(t, "initial", *got)
		}
	})
}

func assertNil[T any](t testing.TB, got *T) {
	t.Helper()

	if got == nil {
		return
	}

	t.Fatalf("assertNil got non-nil:\n%v", got)
}

func assertNotNil[T any](t testing.TB, got *T) {
	t.Helper()

	if got != nil {
		return
	}

	t.Fatalf("assertNotNil got nil %T", got)
}

func assertEq[T any](t testing.TB, want, got T) {
	t.Helper()

	if reflect.DeepEqual(want, got) {
		return
	}

	t.Fatalf(`assertEq failed, got:
%+v
-- want --
%+v
`, got, want)
}
