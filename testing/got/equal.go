package got

import (
	"reflect"
	"testing"
)

// Equal ensures that want and got match.
func Equal(t testing.TB, want, got any) {
	if reflect.DeepEqual(want, got) {
		return
	}

	t.Fatalf(`Equal failed, got:
%+v
-- want --
%+v
`, want, got)
}
