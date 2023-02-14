package gompose

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestReadyOnHttp(t *testing.T) {
	t.Run("marks ready immediately with no options specified", func(t *testing.T) {
		rc := ReadyOnHttp()

		select {
		case err := <-rc:
			assert.NoError(t, err)
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
			assert.NoError(t, err)
			assertServiceIsUp(t)
		case <-time.After(time.Second):
			t.Fatal("was not ready in time")
		}
	})

	t.Run("times out when condition cannot be met", func(t *testing.T) {
		rc := ReadyOnHttp(WithRequest(validRequest(t)), WithTimeout(300*time.Millisecond))

		select {
		case err := <-rc:
			assert.ErrorIs(t, err, ErrWaitTimedOut)
		case <-time.After(400 * time.Millisecond):
			t.Fatal("did not time out in time")
		}
	})

	t.Run("accepts custom request verifiers", func(t *testing.T) {
		testUp(t)
		defer testDown(t)

		bodyIsOk := func(resp *http.Response) (bool, error) {
			bytes, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			defer resp.Body.Close()

			return string(bytes) == "ok\n", nil
		}
		rc := ReadyOnHttp(WithRequest(validRequest(t)), WithResponseVerifier(bodyIsOk))

		select {
		case err := <-rc:
			assert.NoError(t, err)
			assertServiceIsUp(t)
		case <-time.After(time.Second):
			t.Fatal("was not ready in time")
		}
	})
}
