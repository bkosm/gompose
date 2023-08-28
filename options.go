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

func Times(times uint) Option {
	return Option{
		withTimeBasedFunc: func(opt *timeBased) {
			opt.times = times
		},
	}
}

func Timeout(timeout time.Duration) Option {
	return Option{
		withTimeBasedFunc: func(opt *timeBased) {
			opt.timeout = timeout
		},
	}
}

func PollInterval(pollInterval time.Duration) Option {
	return Option{
		withTimeBasedFunc: func(opt *timeBased) {
			opt.pollInterval = pollInterval
		},
	}
}

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

func Wait(channel ReadyOrErrChan) Option {
	return Option{
		withUpFunc: func(opt *up) {
			opt.wait = channel
		},
	}
}

func SignalCallback(callback func(os.Signal)) Option {
	return Option{
		withUpFunc: func(opt *up) {
			opt.onSignal = callback
		},
	}
}

func CustomServices(services ...string) Option {
	return Option{
		withUpFunc: func(opt *up) {
			opt.customServices = services
		},
	}
}
