package gompose

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
	"time"
)

func TestReadyOnStdout(t *testing.T) {
	t.Parallel()

	t.Run("marks ready when a specified phrase occurs in N lines", func(t *testing.T) {
		cmd := exec.Command("echo", "1\n2\n3\n2\n")
		rc := ReadyOnStdout(cmd, WithText("2"), Times(2))

		select {
		case err := <-rc:
			assert.NoError(t, err)
		case <-time.After(time.Second):
			t.Fatal("was not ready in time")
		}
	})

	t.Run("marks ready immediately with no options specified", func(t *testing.T) {
		cmd := exec.Command("pwd")
		rc := ReadyOnStdout(cmd)

		select {
		case err := <-rc:
			assert.NoError(t, err)
		case <-time.After(time.Second):
			t.Fatal("was not ready in time")
		}
	})

	t.Run("fails immediately with command issues", func(t *testing.T) {
		cmd := exec.Command("this-shouldnt-work")
		rc := ReadyOnStdout(cmd)

		select {
		case err := <-rc:
			assert.Error(t, err)
		case <-time.After(time.Second):
			t.Fatal("was not ready in time")
		}
	})

	t.Run("times out after provided duration and returns an error", func(t *testing.T) {
		cmd := exec.Command("pwd")
		rc := ReadyOnStdout(cmd, WithText("nope"), WithTimeout(time.Millisecond))

		select {
		case err := <-rc:
			assert.Error(t, err)
		case <-time.After(2 * time.Millisecond):
			t.Fatal("did not time out in time")
		}
	})

	t.Run("poll interval can be adjusted", func(t *testing.T) {
		c := `
		if [[ ! -e ./pit1 ]]
		then
			touch ./pit1
		else
			rm ./pit1 && echo "ok"
		fi
		`
		cmd := exec.Command("bash", "-c", c)
		rc := ReadyOnStdout(cmd, WithText("ok"), WithPollInterval(time.Millisecond))

		select {
		case err := <-rc:
			assert.NoError(t, err)
		case <-time.After(80 * time.Millisecond): // default is 100ms
			t.Fatal("did not complete in time")
		}
	})
}
