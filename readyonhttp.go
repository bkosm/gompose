package gompose

import "net/http"

// ReadyOnHttp returns a channel that will be closed when the configured http check on the containers is successful.
// The channel will be closed immediately if no options are specified or the request is nil.
// An error will be returned if the request is not nil and the response verifier returns an error,
// and a ErrWaitTimedOut error will be returned if the timeout is reached.
// If the request fails due to a network error, the request will be retried until the timeout is reached.
func ReadyOnHttp(request http.Request, opts ...Option) ReadyOrErrChan {
	var (
		time     timeBased
		verifier responseVerifier
	)
	reduceReadyOnHttpOptions(&time, &verifier, opts)

	readyOrErr := make(chan error)

	go seekOrTimeout(time.timeout, time.pollInterval, readyOrErr, func() (bool, error) {
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

func reduceReadyOnHttpOptions(time *timeBased, verifier *responseVerifier, opts []Option) {
	*time = timeBased{
		times:        1,
		timeout:      DefaultWaitTimeout,
		pollInterval: DefaultPollInterval,
	}
	*verifier = DefaultResponseVerifier

	for _, opt := range opts {
		if fn := opt.withTimeBasedFunc; fn != nil {
			fn(time)
		}
		if fn := opt.withResponseVerifierFunc; fn != nil {
			fn(verifier)
		}
	}
}
