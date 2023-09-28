package gompose

import "net/http"

// ReadyOnHttp returns a channel that will be closed when the configured http check on the containers is successful.
// An error will be returned if the request is successful and the response verifier returns an error,
// and a ErrWaitTimedOut error will be returned if the timeout is reached.
// When it comes to timing defaults, those from ReadyOnStdout apply here too.
// If the request fails due to a network error, the request will be retried until the timeout is reached.
func ReadyOnHttp(request http.Request, opts ...Option) ReadyOrErrChan {
	verifier := responseVerifier(DefaultResponseVerifier)
	options := timeBased{
		times:        1,
		timeout:      DefaultWaitTimeout,
		pollInterval: DefaultPollInterval,
	}
	for _, opt := range opts {
		if fn := opt.withTimeBasedFunc; fn != nil {
			fn(&options)
		}
		if fn := opt.withResponseVerifierFunc; fn != nil {
			fn(&verifier)
		}
	}

	readyOrErr := make(chan error)

	go seekOrTimeout(options.timeout, options.pollInterval, readyOrErr, func() (bool, error) {
		resp, err := http.DefaultClient.Do(&request)
		if err != nil {
			return false, nil // retry it anyways
		}

		ready, err := verifier(resp)
		if err != nil {
			return false, err // can not proceed with waiting, by custom measure
		}

		return ready, nil
	})

	return readyOrErr
}
