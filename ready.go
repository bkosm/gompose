package gompose

import (
	"errors"
	"net/http"
	"time"
)

type (
	// ReadyOrErrChan is a channel that will be closed when the service is ready, or will provide
	// an error with the reason why it failed.
	// An internal reason for failure is a ErrWaitTimedOut, otherwise it will propagate any error that was encountered.
	ReadyOrErrChan <-chan error

	// ReadyOption  is a function that configures how waiting for readiness is performed.
	ReadyOption func(*readyOptions)

	readyOptions struct {
		awaiting         string
		times            uint
		timeout          time.Duration
		pollInterval     time.Duration
		customFile       *string
		request          *http.Request
		responseVerifier func(response *http.Response) (bool, error)
	}
)

// ErrWaitTimedOut error returned in case the specified timeout was reached.
// Without additional configuration it defaults to DefaultWaitTimeout.
var ErrWaitTimedOut = errors.New("gompose: timed out waiting on condition")

const (
	// DefaultPollInterval is the default interval between readiness checks.
	// The interval is performed after the entire checking process has concluded.
	DefaultPollInterval = 100 * time.Millisecond

	// DefaultWaitTimeout is the default timeout for waiting on readiness.
	DefaultWaitTimeout = 10 * time.Minute
)

// DefaultResponseVerifier is the default response verifier used by ReadyOnHttp.
// It verifies that the response status code is 200.
func DefaultResponseVerifier(response *http.Response) (bool, error) {
	return response.StatusCode == http.StatusOK, nil
}

// WithText is a ReadyOption that configures the service to wait until the specified text is found in the output.
func WithText(text string) ReadyOption {
	return func(o *readyOptions) {
		o.awaiting = text
	}
}

// Times specified the number of lines in which the text provided to WithText should be found in the output.
// Default is 1. If 0 is specified, the check will succeed if no error will occur prior.
func Times(n uint) ReadyOption {
	return func(o *readyOptions) {
		o.times = n
	}
}

// WithTimeout specifies how long should the checks be retried before giving up and returning ErrWaitTimedOut.
func WithTimeout(t time.Duration) ReadyOption {
	return func(o *readyOptions) {
		o.timeout = t
	}
}

// WithPollInterval specifies the duration between each check.
func WithPollInterval(t time.Duration) ReadyOption {
	return func(o *readyOptions) {
		o.pollInterval = t
	}
}

// WithRequest specifies the request to be used for ReadyOnHttp.
func WithRequest(req *http.Request) ReadyOption {
	return func(o *readyOptions) {
		o.request = req
	}
}

// WithResponseVerifier specifies the response verifier to be used for ReadyOnHttp.
// DefaultResponseVerifier is used when no other is provided.
func WithResponseVerifier(fn func(response *http.Response) (bool, error)) ReadyOption {
	return func(o *readyOptions) {
		o.responseVerifier = fn
	}
}

// AsReadyOpt converts global GlobalOption's which are useful in the context of awaiting readiness to a ReadyOption.
func AsReadyOpt(fns ...GlobalOption) ReadyOption {
	g := &globalOpts{customFile: nil}
	for _, fn := range fns {
		fn(g)
	}

	return func(o *readyOptions) {
		o.customFile = g.customFile
	}
}

func seekOrTimeout(
	timeout, pollInterval time.Duration,
	readyOrErr chan error,
	seeker func() (bool, error),
) {
	foundOrErr := make(chan error)
	go func() {
		for {
			select {
			case <-foundOrErr:
				return // timeout
			default:
				if found, err := seeker(); err != nil {
					foundOrErr <- err
					return // can not proceed with waiting
				} else if found {
					close(foundOrErr)
					return // found
				} else {
					time.Sleep(pollInterval)
				}
			}
		}
	}()

	select {
	case err := <-foundOrErr:
		if err != nil {
			readyOrErr <- err
		}
		close(readyOrErr)

	case <-time.After(timeout):
		readyOrErr <- ErrWaitTimedOut
		close(readyOrErr)
		close(foundOrErr)
	}
}
