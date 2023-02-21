package gompose

type tt interface {
	Helper()
	Fatal(args ...any)
}

// MustT is a helper function returns a value factory which will eliminate the need to check for errors
// when spawning stubs and fixtures for testing.
// Pass it the instance of *testing.T.
func MustT[T any](t interface{}) func(v T, err error) T {
	tt := t.(tt)
	tt.Helper()

	return func(v T, err error) T {
		if err != nil {
			tt.Fatal(err)
		}
		return v
	}
}
