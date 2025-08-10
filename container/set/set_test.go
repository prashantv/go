package set

import (
	"cmp"
	"maps"
	"reflect"
	"slices"
	"testing"
)

func TestSet_New(t *testing.T) {
	set := New("a", "b")
	assertEq(t, len(set), 2)
	keys := slices.Collect(maps.Keys(set))
	slices.Sort(keys)
	assertEq(t, []string{"a", "b"}, keys)
}

func TestSet_Contains(t *testing.T) {
	set := New("a", "b")
	assertEq(t, true, set.Contains("a"))
	assertEq(t, true, set.Contains("b"))
	assertEq(t, false, set.Contains("c"))
}

func TestSet_ContainsMulti(t *testing.T) {
	set := New("a", "b")

	tests := []struct {
		name            string
		items           []string
		wantContainsAll bool
		wantContainsAny bool
	}{
		{
			name:            "empty",
			items:           nil,
			wantContainsAll: true,
			wantContainsAny: true,
		},
		{
			name:            "single existing",
			items:           arr("a"),
			wantContainsAll: true,
			wantContainsAny: true,
		},
		{
			name:            "single missing",
			items:           arr("missing"),
			wantContainsAll: false,
			wantContainsAny: false,
		},
		{
			name:            "multiple existing",
			items:           arr("b", "a"),
			wantContainsAll: true,
			wantContainsAny: true,
		},
		{
			name:            "some missing",
			items:           arr("b", "missing", "a"),
			wantContainsAll: false,
			wantContainsAny: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("ContainsAll", func(t *testing.T) {
				assertEq(t, tt.wantContainsAll, set.ContainsAll(tt.items))
			})
			t.Run("ContainsAny", func(t *testing.T) {
				assertEq(t, tt.wantContainsAny, set.ContainsAny(tt.items))
			})
		})
	}
}

func TestSet_Copy(t *testing.T) {
	set := New("a")
	set2 := set.Copy()
	assertEq(t, set, set2)
	assertEq(t, true, set.Contains("a"))
	assertEq(t, true, set2.Contains("a"))

	set.Insert("b")
	assertEq(t, true, set.Contains("b"))
	assertEq(t, false, set2.Contains("b"))
}

func TestSet_Delete(t *testing.T) {
	set := New("a", "b")
	assertEq(t, true, set.Contains("a"))
	assertEq(t, true, set.Contains("b"))

	set.Delete("b")
	assertEq(t, true, set.Contains("a"))
	assertEq(t, false, set.Contains("b"))
}

func TestSet_DeleteExists(t *testing.T) {
	set := New("a", "b")
	assertEq(t, true, set.Contains("b"))

	assertEq(t, true, set.DeleteExists("b"))
	assertEq(t, false, set.Contains("b"))

	assertEq(t, false, set.DeleteExists("b"))
	assertEq(t, true, set.Contains("a"))
	assertEq(t, false, set.Contains("b"))
}

func TestSet_Insert(t *testing.T) {
	set := New("a")

	add := []string{"b", "c"}
	for _, item := range add {
		assertEq(t, false, set.Contains(item))
	}

	for _, item := range add {
		set.Insert(item)
		assertEq(t, true, set.Contains(item))

		// Test double insert.
		set.Insert(item)
		assertEq(t, true, set.Contains(item))
	}
}

func TestSet_InsertUnique(t *testing.T) {
	set := New("a")
	assertEq(t, false, set.Contains("b"))

	assertEq(t, true, set.InsertUnique("b"))
	assertEq(t, false, set.InsertUnique("b"))
	assertEq(t, true, set.Contains("b"))
}

func TestSet_InsertSeq(t *testing.T) {
	set := New("a")
	assertEq(t, false, set.Contains("b"))

	add := []string{"a", "b", "c"}
	set.InsertSeq(slices.Values(add))
	for _, item := range add {
		assertEq(t, true, set.Contains(item))
	}
}

func TestSet_Merge(t *testing.T) {
	tests := []struct {
		name          string
		a, b          []string
		wantIntersect []string
		wantUnion     []string
	}{
		{
			name:          "empty",
			a:             nil,
			b:             nil,
			wantIntersect: nil,
			wantUnion:     nil,
		},
		{
			name:          "one empty",
			a:             arr("a", "b"),
			b:             nil,
			wantIntersect: nil,
			wantUnion:     arr("a", "b"),
		},
		{
			name:          "all shared",
			a:             arr("a", "b"),
			b:             arr("a", "b"),
			wantIntersect: arr("a", "b"),
			wantUnion:     arr("a", "b"),
		},
		{
			name:          "no shared",
			a:             arr("a", "b"),
			b:             arr("c", "d"),
			wantIntersect: nil,
			wantUnion:     arr("a", "b", "c", "d"),
		},
		{
			name:          "some shared",
			a:             arr("a", "b"),
			b:             arr("b", "c"),
			wantIntersect: arr("b"),
			wantUnion:     arr("a", "b", "c"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(tt.a...)
			b := New(tt.b...)

			t.Run("Intersect", func(t *testing.T) {
				want := New(tt.wantIntersect...)
				assertEq(t, want, a.Intersect(b))
				assertEq(t, want, b.Intersect(a))
			})

			t.Run("Union", func(t *testing.T) {
				want := New(tt.wantUnion...)
				assertEq(t, want, a.Union(b))
				assertEq(t, want, b.Union(a))
			})
		})
	}
}

func TestSet_Comparisons(t *testing.T) {
	tests := []struct {
		name           string
		a, b           []string
		wantEqual      bool
		wantASubsetB   bool
		wantASupersetB bool
	}{
		{
			name:           "both empty",
			a:              nil,
			b:              nil,
			wantEqual:      true,
			wantASubsetB:   true,
			wantASupersetB: true,
		},
		{
			name:           "same elements",
			a:              arr("a", "b"),
			b:              arr("a", "b"),
			wantEqual:      true,
			wantASubsetB:   true,
			wantASupersetB: true,
		},
		{
			name:         "different elements",
			a:            arr("a", "b"),
			b:            arr("b", "c"),
			wantEqual:    false,
			wantASubsetB: false,
		},
		{
			name:         "subset",
			a:            arr("a"),
			b:            arr("a", "b"),
			wantEqual:    false,
			wantASubsetB: true,
		},
		{
			name:           "superset",
			a:              arr("a", "b"),
			b:              arr("a"),
			wantEqual:      false,
			wantASupersetB: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(tt.a...)
			b := New(tt.b...)

			t.Run("Equal", func(t *testing.T) {
				assertEq(t, tt.wantEqual, a.Equals(b))
				assertEq(t, tt.wantEqual, b.Equals(a))
			})

			t.Run("SubsetOf", func(t *testing.T) {
				assertEq(t, tt.wantASubsetB, a.SubsetOf(b))
			})

			t.Run("SupersetOf", func(t *testing.T) {
				assertEq(t, tt.wantASupersetB, a.SupersetOf(b))
			})
		})
	}
}

func TestSet_All(t *testing.T) {
	tests := []struct {
		name  string
		items []string
	}{
		{
			name:  "empty",
			items: []string{},
		},
		{
			name:  "single item",
			items: []string{"a"},
		},
		{
			name:  "multiple items",
			items: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.items...)
			t.Run("Unordered", func(t *testing.T) {
				got := s.Unordered()

				slices.Sort(got)
				assertEq(t, tt.items, got)
			})

			t.Run("Ordered", func(t *testing.T) {
				got := Ordered(s)
				assertEq(t, tt.items, got)
			})

			t.Run("Iter", func(t *testing.T) {
				got := slices.Collect(s.Iter())
				slices.Sort(got)

				want := tt.items
				if len(want) == 0 {
					// slices.Collect returns nil when there are no items.
					want = nil
				}
				assertEq(t, want, got)
			})
		})
	}
}

func TestSet_Uncomparable(t *testing.T) {
	type S struct {
		V int
	}
	items := []S{{1}, {2}, {3}}
	s := New(items...)
	for _, item := range items {
		assertEq(t, true, s.Contains(item))
	}

	got := s.Unordered()
	slices.SortFunc(got, func(a, b S) int {
		return cmp.Compare(a.V, b.V)
	})
	assertEq(t, items, got)
}

func TestSet_Iter_Break(t *testing.T) {
	s := New("a", "b", "c")

	var got []string
	for item := range s.Iter() {
		got = append(got, item)
		if len(got) == 2 {
			break
		}
	}

	assertEq(t, 2, len(got))
	assertEq(t, true, s.ContainsAll(got))

	gotS := New(got...)
	assertEq(t, false, s.Equals(gotS))
}

func assertEq(t testing.TB, want any, got any) {
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

func arr[T any](vs ...T) []T {
	return vs
}
