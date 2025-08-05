package got

import "testing"

// NoErr is ensures that err is nil.
func NoErr(t testing.TB, err error) {
	t.Helper()

	if err == nil {
		return
	}

	t.Fatalf("NoErr got error: %v", err)
}
