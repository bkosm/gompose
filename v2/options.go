package gompose

import (
	"net/http"
	"os"
	"time"
)

type (
	customFile string
	timeBased  struct {
		times        uint
		timeout      time.Duration
		pollInterval time.Duration
	}
	retry struct {
		times    uint
		interval time.Duration
	}
	responseVerifier func(*http.Response) (bool, error)
	up               struct {
		wait           ReadyOrErrChan
		onSignal       func(os.Signal)
		customServices []string
	}

	// Option is a struct that configures gompose options, shared by all gompose commands.
	Option struct {
		withCustomFileFunc       func(*customFile)
		withTimeBasedFunc        func(*timeBased)
		withResponseVerifierFunc func(*responseVerifier)
		withUpFunc               func(*up)
		withRetryFunc            func(*retry)
	}
)

// CustomFile sets the path of a custom compose file to be used by gompose.
func CustomFile(filepath string) Option {
	return Option{
		withCustomFileFunc: func(opt *customFile) {
			if opt == nil {
				return
			}
			*opt = customFile(filepath)
		},
	}
}

// Times sets the amount of retries that should take place before failing.
func Times(times uint) Option {
	return Option{
		withTimeBasedFunc: func(opt *timeBased) {
			opt.times = times
		},
	}
}

// Timeout sets the duration after which lack of success should propagate failure.
func Timeout(timeout time.Duration) Option {
	return Option{
		withTimeBasedFunc: func(opt *timeBased) {
			opt.timeout = timeout
		},
	}
}

// PollInterval sets the duration that is waited between successive attempts.
func PollInterval(pollInterval time.Duration) Option {
	return Option{
		withTimeBasedFunc: func(opt *timeBased) {
			opt.pollInterval = pollInterval
		},
	}
}

// ResponseVerifier sets the function that verifies a http.Response obtained through ReadyOnHttp.
func ResponseVerifier(verifier func(response *http.Response) (bool, error)) Option {
	return Option{
		withResponseVerifierFunc: func(opt *responseVerifier) {
			if opt == nil {
				return
			}
			*opt = verifier
		},
	}
}

// Wait sets the wait channel for a gompose command.
func Wait(channel ReadyOrErrChan) Option {
	return Option{
		withUpFunc: func(opt *up) {
			opt.wait = channel
		},
	}
}

// SignalCallback sets the callback that is executed whenever a system interrupt happens while the command
// is executing or when the system awaits readiness.
func SignalCallback(callback func(os.Signal)) Option {
	return Option{
		withUpFunc: func(opt *up) {
			opt.onSignal = callback
		},
	}
}

// CustomServices sets a list of services, picked from the default spec or the one provided through CustomFile,
// that should be started with Up.
func CustomServices(services ...string) Option {
	return Option{
		withUpFunc: func(opt *up) {
			opt.customServices = services
		},
	}
}

// RetryCommand will attempt to run the command again in case of failure (e.g. when running Up)
// given amount of times, each after specified interval.
func RetryCommand(times uint, interval time.Duration) Option {
	return Option{
		withRetryFunc: func(opt *retry) {
			opt.times = times
			opt.interval = interval
		},
	}
}
