package gompose

import (
	"errors"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestReadyOnHttp(t *testing.T) {
	t.Run("marks ready immediately with no options specified", func(t *testing.T) {
		t.Parallel()

		rc := ReadyOnHttp()

		select {
		case err := <-rc:
			assertNoError(t, err)
		case <-time.After(time.Second):
			t.Fatal("was not ready in time")
		}
	})

	t.Run("marks ready when default condition (status == 200) is met", func(t *testing.T) {
		testUp(t)
		defer testDown(t)

		rc := ReadyOnHttp(WithRequest(validRequest(t)))

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
			WithRequest(validRequest(t)),
			WithTimeout(2*time.Millisecond),
			WithPollInterval(1*time.Millisecond),
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
		rc := ReadyOnHttp(WithRequest(validRequest(t)), WithResponseVerifier(bodyIsOk))

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

		expected := errors.New("whoops")
		troublemaker := func(resp *http.Response) (bool, error) {
			return false, expected
		}
		rc := ReadyOnHttp(
			WithRequest(validRequest(t)),
			WithResponseVerifier(troublemaker),
			WithPollInterval(time.Second), // to avoid flakiness
		)

		select {
		case err := <-rc:
			assertError(t, err, expected)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("did not fail in time")
		}
	})
}
