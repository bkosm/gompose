package gompose

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"strings"
	"testing"
)

func TestCommand(t *testing.T) {
	t.Run("provides output of an existing command", func(t *testing.T) {
		cmd := *exec.Command("pwd")
		got, err := run(cmd)

		assert.NoError(t, err)
		assert.Greater(t, strings.Index(string(got), "gompose"), 0)
	})

	t.Run("returns error when the command does not exist", func(t *testing.T) {
		cmd := *exec.Command("this-shouldnt-work")
		got, err := run(cmd)

		assert.Error(t, err)
		assert.Empty(t, got)
	})
}
