package gompose

import (
	"fmt"
	"testing"
	"time"
)

func TestDown(t *testing.T) {
	t.Run("down fails if there is no file", func(t *testing.T) {
		t.Parallel()

		err := Down()
		assertError(t, err)
	})

	t.Run("down does not fail if there was no up before", func(t *testing.T) {
		err := Down(customFileOpt)
		assertNoError(t, err)
	})

	t.Run("down cleans up after a successful setup", func(t *testing.T) {
		testUp(t)
		assertEventually(t, serviceIsUp, time.Second, 50*time.Millisecond)

		err := Down(customFileOpt)
		assertNoError(t, err)
		assertServiceIsDown(t)
	})
}

func ExampleDown() {
	err := Down(CustomFile("./testdata/docker-compose.yml"))
	fmt.Print(err)

	// Output:
	// <nil>
}
