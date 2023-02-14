package gompose

import "net/http"

func ReadyOnHttp(fns ...ReadyOption) ReadyOrErrChan {
	opts := &readyOptions{
		pollInterval:     DefaultPollInterval,
		timeout:          DefaultWaitTimeout,
		request:          nil,
		responseVerifier: DefaultResponseVerifier,
	}
	for _, fn := range fns {
		fn(opts)
	}

	readyOrErr := make(chan error)
	if opts.request == nil {
		close(readyOrErr)
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
