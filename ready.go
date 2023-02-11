package gompose

import (
	"errors"
	"time"
)

type ReadyOrErrChan <-chan error

var ErrWaitTimedOut = errors.New("gompose: timed out waiting on condition")

type ReadyOption func(*readyOptions)

type readyOptions struct {
	awaiting     string
	times        uint
	timeout      time.Duration
	pollInterval time.Duration
	customFile   *string
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

func AsReadyOpt(fns ...GomposeOption) ReadyOption {
	g := &gomposeOpts{customFile: nil}
	for _, fn := range fns {
		fn(g)
	}

	return func(o *readyOptions) {
		o.customFile = g.customFile
	}
}
