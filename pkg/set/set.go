package set

import (
	"cmp"
	"slices"
)

// Set is a wrapper for map[T]struct{}.
// Like the underlying map, it is not safe for concurrent use.
type Set[K comparable] map[K]struct{}

func (s Set[K]) Contains(k K) bool {
	_, ok := s[k]
	return ok
}

func (s Set[K]) Insert(k K) {
	s[k] = struct{}{}
}

func (s Set[K]) InsertUnique(k K) bool {
	contains := s.Contains(k)
	if !contains {
		s.Insert(k)
	}
	return contains
}

func (s Set[K]) Delete(k K) {
	delete(s, k)
}

func (s Set[K]) DeleteUnique(k K) bool {
	if !s.Contains(k) {
		return false
	}

	delete(s, k)
	return true
}

func (s Set[K]) ContainsAll(keys []K) bool {
	for _, k := range keys {
		if !s.Contains(k) {
			return false
		}
	}
	return true
}

func (s Set[K]) ContainsAny(keys []K) bool {
	for _, k := range keys {
		if s.Contains(k) {
			return true
		}
	}
	return false
}

func (s Set[K]) Intersects() {

}

// Unordered returns an unordered set of values in the set.
//
// Use `Ordered` when deterministic output is required.
func (s Set[K]) Unordered() []K {
	unordered := make([]K, 0, len(s))
	for k := range s {
		unordered = append(unordered, k)
	}
	return unordered
}

// Ordered returns an ordered set of values in the set.
func Ordered[K cmp.Ordered](s Set[K]) []K {
	ks := s.Unordered()
	slices.Sort(ks)
	return ks
}
