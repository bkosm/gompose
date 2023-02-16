package gompose

import "testing"

func MustT[T any](t *testing.T) func(v T, err error) T {
	t.Helper()

	return func(v T, err error) T {
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
}
