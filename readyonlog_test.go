package gompose

import (
	"fmt"
	"testing"
	"time"
)

func TestReadyOnLog(t *testing.T) {
	t.Run("ready on log pushes compose logs into ReadyOnStdout", func(t *testing.T) {
		testUp(t)
		defer testDown(t)

		rc := ReadyOnLog(expectedLogLine, customFileOpt)

		select {
		case err := <-rc:
			assertNoError(t, err)
		case <-time.After(2 * time.Minute):
			t.Fatal("time out waiting on compose (might be pulling the image)")
		}
	})
}

func ExampleReadyOnLog() {
	_ = Up(customFileOpt)
	ch := ReadyOnLog(expectedLogLine, customFileOpt)

	<-ch
	fmt.Println("the service is up now")
	// Output:
	// the service is up now

	_ = Down(customFileOpt)
}
