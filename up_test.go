package gompose

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestUp(t *testing.T) {
	t.Run("up fails if there is no file", func(t *testing.T) {
		err := Up()
		assert.Error(t, err)
	})

	t.Run("up works with no arguments", func(t *testing.T) {
		goBack := goIntoTestDataDir(t)
		defer func() {
			goBack()
			testDown(t)
		}()

		err := Up()
		assert.NoError(t, err)
		assertServiceIsUp(t)
	})

	t.Run("up works with options", func(t *testing.T) {
		defer testDown(t)

		err := Up(
			WithWait(ReadyOnLog(WithText(expectedLogLine), AsReadyOpt(customFileOpt))),
			WithCustomServices(customServiceName),
			AsUpOpt(customFileOpt),
		)
		assert.NoError(t, err)
		assertServiceIsUp(t)
	})

	t.Run("intercepts os signals", func(t *testing.T) {
		defer testDown(t)

		c := 0
		callback := func(s os.Signal) {
			if s == os.Interrupt {
				c += 1
			}
		}

		err := Up(
			WithWait(ReadyOnLog(WithText(expectedLogLine), AsReadyOpt(customFileOpt))),
			WithSignalCallback(callback),
			AsUpOpt(customFileOpt),
		)
		assert.NoError(t, err)
		assertServiceIsUp(t)

		doSignal(t, syscall.SIGINT)

		wasCalled := func() bool { return c == 1 }
		assert.Eventually(t, wasCalled, 5*time.Second, 100*time.Millisecond)
	})

	t.Run("propagates wait channel errors", func(t *testing.T) {
		defer testDown(t)
		c, done := make(chan error), make(chan any)
		e := errors.New("whoops")

		go func() {
			err := Up(WithWait(c), AsUpOpt(customFileOpt))
			assert.ErrorIs(t, err, e)
			close(done)
		}()

		c <- e
		<-done
	})
}
