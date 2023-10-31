package gompose

import (
	"io"
	"net/http"
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
			Wait(
				ReadyOnLog(expectedLogLine),
			),
			SignalCallback(func(_ os.Signal) {
				_ = Down()
			}),
			RetryCommand(3, time.Second),
		)
		assertNoError(t, err)
		assertServiceIsUp(t)
	}

	teardown := func() {
		err := Down()
		assertNoError(t, err)
		assertServiceIsDown(t)
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
		teardown()
	})

	t.Run("allows for waiting on healthy http", func(t *testing.T) {
		req := validRequest(t)

		err := Up(Wait(ReadyOnHttp(req)))
		assertNoError(t, err)
		assertServiceIsUp(t)

		teardown()
	})

	t.Run("allows customising the response verifiers", func(t *testing.T) {
		req := validRequest(t)

		verifier := ResponseVerifier(func(res *http.Response) (bool, error) {
			b, err := io.ReadAll(res.Body)
			if err != nil {
				return false, err
			}

			return string(b) == "ok\n", nil
		})

		err := Up(Wait(ReadyOnHttp(req, verifier, Timeout(3*time.Second))))
		assertNoError(t, err)
		assertServiceIsUp(t)

		teardown()
	})

	t.Run("skips Down when an environment flag is present", func(t *testing.T) {
		err := os.Setenv(SkipEnv, "DOWN,IGNORE")
		assertNoError(t, err)

		setup()

		err = Down()
		assertNoError(t, err)
		assertServiceIsUp(t)

		err = os.Unsetenv(SkipEnv)
		assertNoError(t, err)
		teardown()
	})
}
