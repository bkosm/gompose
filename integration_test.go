package gompose

import (
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
		assertNoError(t, err)
		assertServiceIsUp(t)
	}

	t.Run("sets up the services", func(t *testing.T) {
		setup()
	})

	t.Run("cleans up on system signals", func(t *testing.T) {
		doSignal(t, syscall.SIGINT)
		assertEventually(t, serviceIsDown, 5*time.Second, 100*time.Millisecond)
	})

	t.Run("sets up again after a forced exit", func(t *testing.T) {
		setup()
	})

	t.Run("cleans up on direct request", func(t *testing.T) {
		err := Down()
		assertNoError(t, err)
		assertServiceIsDown(t)
	})

	t.Run("allows for waiting on healthy http", func(t *testing.T) {
		req := validRequest(t)

		err := Up(WithWait(ReadyOnHttp(WithRequest(req))))
		assertNoError(t, err)
		assertServiceIsUp(t)
		assertNoError(t, Down())
	})
}
