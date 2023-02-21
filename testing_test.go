package gompose

import (
	"fmt"
	"net/http"
	"testing"
)

type ft struct {
	helperCalls uint
	fatalCalls  uint
}

func (f *ft) Helper() {
	f.helperCalls++
}

func (f *ft) Fatal(_ ...any) {
	f.fatalCalls++
}

func TestMustTWorks(t *testing.T) {
	t.Parallel()

	t.Run("records errors", func(t *testing.T) {
		t.Parallel()

		ft := &ft{}
		fn := MustT[*http.Request](ft)

		val := fn(http.NewRequest("", "\n", nil))

		if ft.helperCalls != 1 {
			t.Fatal("helper was not called")
		}
		if ft.fatalCalls != 1 {
			t.Fatal("did not record the error")
		}
		if val != nil {
			t.Fatal("did not return nil")
		}
	})

	t.Run("returns the value", func(t *testing.T) {
		t.Parallel()

		ft := &ft{}
		fn := MustT[*http.Request](ft)

		val := fn(http.NewRequest("", "", nil))

		if ft.helperCalls != 1 {
			t.Fatal("helper was not called")
		}
		if ft.fatalCalls != 0 {
			t.Fatal("recorded an error")
		}
		if val == nil {
			t.Fatal("returned nil")
		}
	})

	t.Run("panics when the required methods are not present on the argument", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if err := recover(); err == nil {
				t.Fatal("did not panic")
			}
		}()

		fn := MustT[*http.Request](interface{}(nil))
		_ = fn(http.NewRequest("", "", nil))
	})
}

func ExampleMustT() {
	t := &ft{}
	fn := MustT[*http.Request](t)

	rq := fn(http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d", containerPort), nil))
	fmt.Println(rq.URL.Scheme)
	// Output:
	// http
}
