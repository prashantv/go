package got

import (
	"fmt"
	"reflect"
	"testing"
)

// EmptyV is a generic wrapper for Empty.
func EmptyV[V any](t testing.TB, v V) {
	Empty(t, v)
}

// Empty ensures that the given value is empty.
// Empty differs from 0 as it treats non-nil contains as empty if their length is 0.
// Arrays are empty if every value is empty.
// For struct, it checks if every field is empty.
func Empty(t testing.TB, v any) {
	t.Helper()

	if v == nil {
		return
	}

	err := empty(reflect.ValueOf(v))
	if err == nil {
		return
	}

	t.Fatalf("Empty: %T: %v:\n%+v\n", v, err, v)
}

func empty(rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Chan, reflect.Slice, reflect.Map:
		if rv.IsNil() {
			return nil
		}

		if rv.Len() == 0 {
			return nil
		}

		return fmt.Errorf("container not empty, %v elements", rv.Len())
	case reflect.Array:
		// TODO: check each elem
		panic("unimplemented")
	case reflect.Func:
		if rv.IsNil() {
			return nil
		}

		return fmt.Errorf("funt non-nil")
	case reflect.Ptr:
		if rv.IsNil() {
			return nil
		}

		if err := empty(rv.Elem()); err != nil {
			return fmt.Errorf(".*: %w", err)
		}

		return nil
	case reflect.Struct:
		rt := rv.Type()
		for i := range rt.NumField() {
			if err := empty(rv.Field(i)); err != nil {
				return fmt.Errorf("%v: %w", rt.Field(i).Name, err)
			}
		}
		return nil
	default:
		if rv.IsZero() {
			return nil
		}

		return fmt.Errorf("non-zero value")
	}
}
