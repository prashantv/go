package retaintest

import (
	"strings"
	"testing"

	"github.com/goccy/go-yaml"

	"go.prashantv.com/testing/got"
)

type S struct {
	Retain map[any]any `yaml:",inline"`

	Name string `json:"name,omitempty"`
}

func TestRetain_RoundTrip(t *testing.T) {
	// Tests based on from https://github.com/prashantv/pkg/tree/main/jsonobj
	tests := []struct {
		name   string
		json   string
		verify func(testing.TB, S)

		update         func(*S)
		wantUpdateJSON string
	}{
		{
			name:   "empty",
			json:   `{}`,
			verify: got.EmptyV[S],
		},
		{
			name: "only known",
			json: `name: foo`,
			verify: func(t testing.TB, s S) {
				got.Equal(t, "foo", s.Name)
			},

			update: func(s *S) {
				s.Name = "new name"
			},
			wantUpdateJSON: `name: new name`,
		},
		{
			name: "known and unknown",
			json: `
list:
- 1
- 2
- 3
num: 1
obj:
  k: v
str: foo
name: foo
`,
			verify: func(t testing.TB, s S) {
				got.Equal(t, "foo", s.Name)
			},

			update: func(s *S) {
				s.Name = "new name"
			},
			wantUpdateJSON: `
list:
- 1
- 2
- 3
num: 1
obj:
  k: v
str: foo
name: new name
`,
		},
	}

	multiline := func(s string) string {
		return strings.TrimSpace(s) + "\n"
	}

	for _, tt := range tests {
		tt.json = strings.TrimSpace(tt.json)
		t.Run(tt.name, func(t *testing.T) {
			var s S
			err := yaml.Unmarshal([]byte(tt.json), &s)
			got.NoErr(t, err)

			t.Run("verify Unmarshal", func(t *testing.T) {
				tt.verify(t, s)

				got.Equal(t, multiline(tt.json), mustMarshal(t, &s))
			})

			if tt.update != nil {
				t.Run("verify Marshal after mutate", func(t *testing.T) {
					tt.update(&s)
					got.Equal(t, multiline(tt.wantUpdateJSON), mustMarshal(t, &s))
				})
			}
		})
	}
}

func mustMarshal(t testing.TB, obj any) string {
	v, err := yaml.Marshal(obj)
	got.NoErr(t, err)
	return string(v)
}
