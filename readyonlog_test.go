package gompose

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReadyOnLog(t *testing.T) {
	t.Run("ready on log pushes compose logs into ReadyOnStdout", func(t *testing.T) {
		testUp(t)
		defer testDown(t)

		customFileOpt := WithCustomFile("./testdata/docker-compose.yml")
		rc := ReadyOnLog(Text(expectedLine), AsReadyOpt(customFileOpt))

		select {
		case err := <-rc:
			assert.NoError(t, err)
		case <-time.After(5 * time.Minute):
			t.Fatal("time out waiting on compose (might be pulling the image)")
		}
	})
}