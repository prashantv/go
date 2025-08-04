package xslices

import (
	"reflect"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []string
	}{
		{
			name: "nil",
			in:   nil,
			want: nil,
		},
		{
			name: "empty slice",
			in:   []int{},
			want: []string{},
		},
		{
			name: "single element",
			in:   []int{1},
			want: []string{"1"},
		},
		{
			name: "multiple elements",
			in:   []int{1, 2, 3},
			want: []string{"1", "2", "3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Map(tt.in, strconv.Itoa)
			assertEq(t, tt.want, got)
		})
	}
}

func BenchmarkMap(b *testing.B) {
	xs := make([]int, 10000)
	for b.Loop() {
		Map(xs, func(x int) int {
			return x
		})
	}
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
