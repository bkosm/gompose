package gompose

import (
	"errors"
	"testing"
	"time"
)

func Test_doRetry(t *testing.T) {
	t.Run("should return immediately in case of no error", func(t *testing.T) {
		start := time.Now()
		defer func() {
			if time.Now().Sub(start) > time.Second {
				t.Fail()
			}
		}()

		err := doRetry(1, time.Second, func() error {
			return nil
		})
		assertNoError(t, err)
	})

	t.Run("should return error after all attempts", func(t *testing.T) {
		start := time.Now()
		defer func() {
			if time.Now().Sub(start) < time.Millisecond*10 {
				t.Fail()
			}
		}()

		err := doRetry(3, time.Millisecond*10, func() error {
			return errors.New("boom")
		})
		assertError(t, err)
	})

	t.Run("should return success after a retry if the operation passed", func(t *testing.T) {
		start := time.Now()
		defer func() {
			if time.Now().Sub(start) < time.Millisecond*10 {
				t.Fail()
			}
		}()

		var (
			i        = 0
			innerErr = errors.New("boom")
		)
		err := doRetry(3, time.Millisecond*10, func() error {
			if i == 2 {
				return nil
			}
			i++
			return innerErr
		})
		assertNoError(t, err)
	})
}
