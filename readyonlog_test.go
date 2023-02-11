package gompose

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReadyOnLog(t *testing.T) {
	customFileOpt := WithCustomFile("./testdata/docker-compose.yml")

	t.Run("ready on log pushes compose logs into ReadyOnStdout", func(t *testing.T) {
		testUp(t)
		defer testDown(t)

		rc := ReadyOnLog(Text(expectedLine), AsReadyOpt(customFileOpt))

		select {
		case err := <-rc:
			assert.NoError(t, err)
		case <-time.After(5 * time.Minute):
			t.Fatal("time out waiting on compose (might be pulling the image)")
		}
	})

	t.Run("marks ready immediately with no options specified", func(t *testing.T) {
		testUp(t)
		defer testDown(t)

		// it needs a lot extra code for no custom file stuff
		rc := ReadyOnLog(AsReadyOpt(customFileOpt))

		select {
		case err := <-rc:
			assert.NoError(t, err)
		case <-time.After(time.Second):
			t.Fatal("was not ready in time")
		}
	})
}
