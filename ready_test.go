package gompose

import (
	"os/exec"
	"testing"
	"time"
)

func TestReadyOnStdout(t *testing.T) {
	t.Parallel()

	t.Run("marks ready when a specified phrase occurs in N lines", func(t *testing.T) {
		cmd := exec.Command("echo", "1\n2\n3\n2\n")

		rc := ReadyOnStdout(cmd, AwaitingText("2"), Times(2))

		select {
		case <-rc:
		case <-time.After(time.Second):
			t.Fatal("was not ready in time")
		}
	})
}
