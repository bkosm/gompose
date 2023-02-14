package gompose

import (
	"errors"
	"net/http"
	"time"
)

type ReadyOrErrChan <-chan error
type ReadyOption func(*readyOptions)

var ErrWaitTimedOut = errors.New("gompose: timed out waiting on condition")

const (
	DefaultPollInterval = 100 * time.Millisecond
	DefaultWaitTimeout  = 10 * time.Minute
)

type readyOptions struct {
	awaiting     string
	times        uint
	timeout      time.Duration
	pollInterval time.Duration
	customFile   *string
	request      *http.Request
}

func WithText(text string) ReadyOption {
	return func(o *readyOptions) {
		o.awaiting = text
	}
}

func Times(n uint) ReadyOption {
	return func(o *readyOptions) {
		o.times = n
	}
}

func WithTimeout(t time.Duration) ReadyOption {
	return func(o *readyOptions) {
		o.timeout = t
	}
}

func WithPollInterval(t time.Duration) ReadyOption {
	return func(o *readyOptions) {
		o.pollInterval = t
	}
}

func WithRequest(req *http.Request) ReadyOption {
	return func(o *readyOptions) {
		o.request = req
	}
}

func AsReadyOpt(fns ...GomposeOption) ReadyOption {
	g := &gomposeOpts{customFile: nil}
	for _, fn := range fns {
		fn(g)
	}

	return func(o *readyOptions) {
		o.customFile = g.customFile
	}
}
