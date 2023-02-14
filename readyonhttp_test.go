package gompose

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReadyOnHttp(t *testing.T) {
	t.Run("marks ready immediately with no options specified", func(t *testing.T) {
		rc := ReadyOnHttp()

		select {
		case err := <-rc:
			assert.NoError(t, err)
		case <-time.After(time.Second):
			t.Fatal("was not ready in time")
		}
	})
}
