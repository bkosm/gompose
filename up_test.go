package gompose

import (
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestUp(t *testing.T) {
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

	customFileOpt := WithCustomFile("./testdata/docker-compose.yml")

	t.Run("up works with options", func(t *testing.T) {
		defer testDown(t)

		err := Up(
			WaitFor(ReadyOnLog(Text(expectedLine), AsReadyOpt(customFileOpt))),
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
			WaitFor(ReadyOnLog(Text(expectedLine), AsReadyOpt(customFileOpt))),
			OnSignal(callback),
			AsUpOpt(customFileOpt),
		)
		assert.NoError(t, err)
		assertServiceIsUp(t)

		doSignal(t, syscall.SIGINT)
		time.Sleep(200 * time.Millisecond)
		assert.Equal(t, 1, c)
	})
}
