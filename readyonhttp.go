package gompose

import "net/http"

// ReadyOnHttp returns a channel that will be closed when the configured http check on the containers is successful.
// The channel will be closed immediately if no options are specified or the request is nil.
// An error will be returned if the request is not nil and the response verifier returns an error,
// and a ErrWaitTimedOut error will be returned if the timeout is reached.
// If the request fails due to a network error, the request will be retried until the timeout is reached.
func ReadyOnHttp(fns ...ReadyOption) ReadyOrErrChan {
	opts := readyOptions{
		pollInterval:     DefaultPollInterval,
		timeout:          DefaultWaitTimeout,
		request:          nil,
		responseVerifier: DefaultResponseVerifier,
	}
	for _, fn := range fns {
		fn(&opts)
	}

	readyOrErr := make(chan error)
	if opts.request == nil {
		close(readyOrErr)
		return readyOrErr // deliberate configuration, ready immediately
	}

	go seekOrTimeout(opts.timeout, opts.pollInterval, readyOrErr, func() (bool, error) {
		resp, err := http.DefaultClient.Do(opts.request)
		if err != nil {
			return false, nil // retry it anyways
		}

		ready, err := opts.responseVerifier(resp)
		if err != nil {
			return false, err // can not proceed with waiting, by custom measure
		}

		return ready, nil
	})

	return readyOrErr
}
