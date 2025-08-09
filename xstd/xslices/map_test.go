package xslices

import (
	"context"
	"reflect"
	"strconv"
	"strings"
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

func TestMapErr(t *testing.T) {
	tests := []struct {
		name    string
		in      []string
		want    []int
		wantErr string
	}{
		{
			name: "nil",
			in:   nil,
			want: nil,
		},
		{
			name: "empty",
			in:   []string{},
			want: []int{},
		},
		{
			name: "single element",
			in:   []string{"1"},
			want: []int{1},
		},
		{
			name: "multiple elmeents",
			in:   []string{"1", "2", "3"},
			want: []int{1, 2, 3},
		},
		{
			name:    "err",
			in:      []string{"1", "2", "err", "3"},
			wantErr: "index 2: strconv.Atoi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapErr(tt.in, strconv.Atoi)
			assertErr(t, tt.wantErr, err)
			assertEq(t, tt.want, got)
		})
	}
}

func TestMapCtx(t *testing.T) {
	canceled := func() context.Context {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		return ctx
	}()

	ignoreCtx := func(_ context.Context, s string) (int, error) {
		return strconv.Atoi(s)
	}
	respectCtx := func(ctx context.Context, s string) (int, error) {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		return ignoreCtx(ctx, s)
	}

	tests := []struct {
		name    string
		fn      []func(context.Context, string) (int, error)
		ctx     []context.Context
		in      []string
		want    []int
		wantErr string
	}{
		{
			name: "nil",
			fn:   arr(ignoreCtx, respectCtx),
			ctx:  arr(context.Background(), canceled),
			in:   nil,
			want: nil,
		},
		{
			name: "empty",
			fn:   arr(ignoreCtx, respectCtx),
			ctx:  arr(context.Background(), canceled),
			in:   []string{},
			want: []int{},
		},
		{
			name: "non-empty",
			fn:   arr(ignoreCtx, respectCtx),
			ctx:  arr(context.Background()),
			in:   []string{"1", "2", "3"},
			want: []int{1, 2, 3},
		},
		{
			name: "ignore context error",
			fn:   arr(ignoreCtx),
			ctx:  arr(context.Background(), canceled),
			in:   []string{"1", "2", "3"},
			want: []int{1, 2, 3},
		},
		{
			name:    "context error",
			fn:      arr(respectCtx),
			ctx:     arr(canceled),
			in:      []string{"1", "2", "3"},
			wantErr: "index 0: " + context.Canceled.Error(),
		},
		{
			name:    "fn error",
			fn:      arr(ignoreCtx),
			ctx:     arr(context.Background(), canceled),
			in:      []string{"1", "err", "3"},
			wantErr: "index 1: strconv.Atoi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, fn := range tt.fn {
				for _, ctx := range tt.ctx {
					got, err := MapCtx(ctx, tt.in, fn)
					assertErr(t, tt.wantErr, err)
					assertEq(t, tt.want, got)
				}
			}
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

func assertErr(t testing.TB, wantErr string, err error) {
	t.Helper()

	if wantErr == "" {
		if err != nil {
			t.Fatalf("assertErr failed, want no error, got:\n%v", err)
		}
		return
	}

	if err == nil {
		t.Fatalf("assertErr failed, wanted error, got nil. wante:\n%v", wantErr)
	}

	if !strings.Contains(err.Error(), wantErr) {
		t.Fatalf(`assertErr failed, got unexpected error:
%v
-- want (contains) --
%v
`, err, wantErr)
	}
}

func arr[T any](vs ...T) []T {
	return vs
}
