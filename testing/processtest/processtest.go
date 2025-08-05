package processtest

import (
	"os"
	"os/exec"
	"testing"
)

// Run should be used at the top-level of a test to run a given child process.
// It returns whether the caller should return.
func Run(t testing.TB, arg string, subProc func(t testing.TB, arg string)) bool {
	k := "TEST_SUBPROCESS_" + t.Name()
	if v := os.Getenv(k); v != "" {
		// We are the child process, process the data.
		subProc(t, v)

		// Make sure test stops at this point.
		return true
	}

	cmd := exec.Command(os.Args[0], "-test.run", t.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	os.Setenv(k, arg)
	return false
}
