package set

import (
	"cmp"
	"iter"
	"slices"
)

// Set implements set operations using a `map[T]struct{}`.
// It is not safe for concurrent use.
type Set[T comparable] map[T]struct{}

// New creates a set with items.
func New[T comparable](items ...T) Set[T] {
	s := make(Set[T], len(items))
	for _, item := range items {
		s.Insert(item)
	}
	return s
}

// Contains returns if the set contains the specified item.
func (s Set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

// ContainsAll returns if all the items exist in the set.
func (s Set[T]) ContainsAll(items []T) bool {
	for _, item := range items {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// ContainsAny returns true if any of the items exist in the set.
// If no items are specified, it returns false.
func (s Set[T]) ContainsAny(items []T) bool {
	if len(items) == 0 {
		return true
	}

	for _, k := range items {
		if s.Contains(k) {
			return true
		}
	}
	return false
}

// Copy returns a new set with the same items.
func (s Set[T]) Copy() Set[T] {
	clone := make(Set[T], len(s))
	for item := range s {
		clone.Insert(item)
	}
	return clone
}

// Insert inserts the item into the set, overwriting any existing items.
func (s Set[T]) Insert(item T) {
	s[item] = struct{}{}
}

// InsertUnique inserts the item into the set if the item is not already in the set.
// It returns true if the item did not previously exist, and was inserted.
func (s Set[T]) InsertUnique(item T) bool {
	contains := s.Contains(item)
	if !contains {
		s.Insert(item)
		return true
	}
	return false
}

// InsertSeq inserts all values from seq into the set, overwriting any existing items.
func (s Set[T]) InsertSeq(seq iter.Seq[T]) {
	for item := range seq {
		s.Insert(item)
	}
}

// Intersect returns a set that only contains items that are in both sets.
func (s Set[T]) Intersect(other Set[T]) Set[T] {
	intersect := make(Set[T])
	for item := range s {
		if other.Contains(item) {
			intersect.Insert(item)
		}
	}
	return intersect
}

// Delete deletes the item from the set.
func (s Set[T]) Delete(item T) {
	delete(s, item)
}

// DeleteExists deletes the item from the set if it exists.
// It returns true if the item was deleted.
func (s Set[T]) DeleteExists(item T) bool {
	if !s.Contains(item) {
		return false
	}

	delete(s, item)
	return true
}

// Equals returns if the two sets are equal.
func (s Set[T]) Equals(other Set[T]) bool {
	if len(s) != len(other) {
		return false
	}
	return s.SubsetOf(other)
}

// SubsetOf returns if other contains all elements in s.
func (s Set[T]) SubsetOf(other Set[T]) bool {
	for item := range s {
		if !other.Contains(item) {
			return false
		}
	}
	return true
}

// SupersetOf returns if s contains all elements in other.
func (s Set[T]) SupersetOf(other Set[T]) bool {
	return other.SubsetOf(s)
}

// Unordered returns an unordered set of values in the set.
// Since it relies on Go map iteration order, the order of the values is non-deterministic.
//
// Use `Ordered` when deterministic output is required.
func (s Set[K]) Unordered() []K {
	unordered := make([]K, 0, len(s))
	for k := range s {
		unordered = append(unordered, k)
	}
	return unordered
}

// Union returns a set with elements from both sets.
func (s Set[T]) Union(other Set[T]) Set[T] {
	union := make(Set[T], max(len(s), len(other)))
	for item := range s {
		union.Insert(item)
	}
	for item := range other {
		union.Insert(item)
	}
	return union
}

// Iter returns an iterator over all items in the set.
func (s Set[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for item := range s {
			if !yield(item) {
				return
			}
		}
	}
}

// Ordered returns an ordered set of values in the set.
func Ordered[T cmp.Ordered](s Set[T]) []T {
	ks := s.Unordered()
	slices.Sort(ks)
	return ks
}
