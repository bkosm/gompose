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
)

// ErrWaitTimedOut error returned in case the specified timeout was reached.
// Without additional configuration it defaults to DefaultWaitTimeout.
var ErrWaitTimedOut = errors.New("gompose: timed out waiting on condition")

const (
	// DefaultPollInterval is the default interval between readiness checks.
	// The interval is performed after the entire checking process has concluded.
	DefaultPollInterval = 100 * time.Millisecond

	// DefaultWaitTimeout is the default timeout for waiting on readiness.
	DefaultWaitTimeout = 5 * time.Minute

	// SkipEnv is the key of the environment variable that is used for flagging.
	SkipEnv = "GOMPOSE_SKIP"

	// SkipDown is the flag used for skipping invoking compose Down without altering source code.
	SkipDown = "down"
)

// DefaultResponseVerifier is the default response verifier used by ReadyOnHttp.
// It verifies that the response status code is 200.
func DefaultResponseVerifier(response *http.Response) (bool, error) {
	return response.StatusCode == http.StatusOK, nil
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

func reduceCustomFileOptions(opts []Option) customFile {
	var file customFile

	for _, opt := range opts {
		if fn := opt.withCustomFileFunc; fn != nil {
			fn(&file)
		}
	}

	return file
}
