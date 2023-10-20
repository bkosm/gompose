package gompose

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	t.Run("provides output of an existing command", func(t *testing.T) {
		t.Parallel()

		cmd := *exec.Command("pwd")
		got, err := run(cmd)

		assertNoError(t, err)
		if !strings.Contains(string(got), "gompose") {
			t.Fatal("expected output to contain 'gompose', got", string(got))
		}
	})

	t.Run("returns error when the command does not exist", func(t *testing.T) {
		t.Parallel()

		cmd := *exec.Command("this-shouldnt-work")
		got, err := run(cmd)

		assertError(t, err)
		if got != "" {
			t.Fatal("expected empty output, got", string(got))
		}
	})
}
