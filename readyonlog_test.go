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

		rc := ReadyOnLog(WithText(expectedLogLine), AsReadyOpt(customFileOpt))

		select {
		case err := <-rc:
			assertNoError(t, err)
		case <-time.After(2 * time.Minute):
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
			assertNoError(t, err)
		case <-time.After(time.Second):
			t.Fatal("was not ready in time")
		}
	})
}

func ExampleReadyOnLog() {
	_ = Up(AsUpOpt(customFileOpt))
	ch := ReadyOnLog(WithText(expectedLogLine), AsReadyOpt(customFileOpt))

	<-ch
	fmt.Println("the service is up now")
	// Output:
	// the service is up now

	_ = Down(AsDownOpt(customFileOpt))
}
