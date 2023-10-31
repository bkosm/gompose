package gompose

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestReadyOnHttp(t *testing.T) {
	t.Run("marks ready when default condition (status == 200) is met", func(t *testing.T) {
		testUp(t)
		defer testDown(t)

		rc := ReadyOnHttp(validRequest(t))

		select {
		case err := <-rc:
			assertNoError(t, err)
			assertServiceIsUp(t)
		case <-time.After(time.Second):
			t.Fatal("was not ready in time")
		}
	})

	t.Run("times out when condition cannot be met", func(t *testing.T) {
		rc := ReadyOnHttp(
			validRequest(t),
			Timeout(2*time.Millisecond),
			PollInterval(1*time.Millisecond),
		)

		select {
		case err := <-rc:
			assertError(t, err, ErrWaitTimedOut)
		case <-time.After(4 * time.Millisecond):
			t.Fatal("did not time out in time")
		}
	})

	t.Run("accepts custom request verifiers", func(t *testing.T) {
		testUp(t)
		defer testDown(t)

		bodyIsOk := func(resp *http.Response) (bool, error) {
			bytes, err := io.ReadAll(resp.Body)
			assertNoError(t, err)

			defer resp.Body.Close()

			return string(bytes) == "ok\n", nil
		}
		rc := ReadyOnHttp(validRequest(t), ResponseVerifier(bodyIsOk))

		select {
		case err := <-rc:
			assertNoError(t, err)
			assertServiceIsUp(t)
		case <-time.After(time.Second):
			t.Fatal("was not ready in time")
		}
	})

	t.Run("custom verifier can return immediate error", func(t *testing.T) {
		testUp(t)
		defer testDown(t)

		assertEventually(t, serviceIsUp, time.Second, 100*time.Millisecond)

		expected := errors.New("whoops")
		troublemaker := ResponseVerifier(func(_ *http.Response) (bool, error) {
			return false, expected
		})

		rc := ReadyOnHttp(
			validRequest(t),
			troublemaker,
			PollInterval(time.Second), // to avoid flakiness
		)

		select {
		case err := <-rc:
			assertError(t, err, expected)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("did not fail in time")
		}
	})
}

func ExampleReadyOnHttp() {
	_ = Up(CustomFile("./testdata/docker-compose.yml"))
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d", containerPort), nil)
	ch := ReadyOnHttp(*request)

	<-ch
	fmt.Println("the service is up now")
	// Output:
	// the service is up now

	_ = Down(customFileOpt)
}
