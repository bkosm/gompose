package gompose

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestUp(t *testing.T) {
	t.Run("up fails if there is no file", func(t *testing.T) {
		t.Parallel()

		err := Up()
		assertError(t, err)
	})

	t.Run("up works with no arguments", func(t *testing.T) {
		goBack := goIntoTestDataDir(t)
		defer func() {
			goBack()
			testDown(t)
		}()

		err := Up()
		assertNoError(t, err)
		assertEventually(t, serviceIsUp, time.Second, 100*time.Millisecond)
	})

	t.Run("up works with options", func(t *testing.T) {
		defer testDown(t)

		err := Up(
			Wait(ReadyOnLog(expectedLogLine, customFileOpt)),
			CustomServices(customServiceName),
			customFileOpt,
		)
		assertNoError(t, err)
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
			Wait(ReadyOnLog(expectedLogLine, customFileOpt)),
			SignalCallback(callback),
			customFileOpt,
		)
		assertNoError(t, err)
		assertServiceIsUp(t)

		doSignal(t, syscall.SIGINT)

		wasCalled := func() bool { return c == 1 }
		assertEventually(t, wasCalled, 5*time.Second, 100*time.Millisecond)
	})

	t.Run("propagates wait channel errors", func(t *testing.T) {
		defer testDown(t)
		c, done := make(chan error), make(chan any)
		expected := errors.New("whoops")

		go func() {
			err := Up(Wait(c), customFileOpt)
			assertError(t, err, expected)
			close(done)
		}()

		c <- expected
		<-done
	})
}

func ExampleUp() {
	_ = Up(
		Wait(ReadyOnLog(expectedLogLine, CustomFile("./testdata/docker-compose.yml"))),
		CustomServices(customServiceName),
		customFileOpt,
	)

	fmt.Println("the containers are ready to go!")
	// Output:
	// the containers are ready to go!

	_ = Down(customFileOpt)
}
