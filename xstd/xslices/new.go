package xslices

// New creates a slice for the given values.
// This is helpful to avoid noise when
func New[T any](vs ...T) []T {
	return vs
}
