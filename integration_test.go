package gompose

import (
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	goBack := goIntoTestDataDir(t)
	defer goBack()

	setup := func() {
		err := Up(
			WithWait(
				ReadyOnLog(WithText(expectedLogLine)),
			),
			WithSignalCallback(func(_ os.Signal) {
				_ = Down()
			}),
		)
		assert.NoError(t, err)
		assertServiceIsUp(t)
	}

	t.Run("sets up the services", func(t *testing.T) {
		setup()
	})

	t.Run("cleans up on system signals", func(t *testing.T) {
		doSignal(t, syscall.SIGTERM)
		time.Sleep(300 * time.Millisecond)
		assertServiceIsDown(t)
	})

	t.Run("sets up again after a forced exit", func(t *testing.T) {
		setup()
	})

	t.Run("cleans up on direct request", func(t *testing.T) {
		err := Down()
		assert.NoError(t, err)
		assertServiceIsDown(t)
	})
}
