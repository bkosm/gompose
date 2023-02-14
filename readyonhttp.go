package gompose

func ReadyOnHttp(fns ...ReadyOption) ReadyOrErrChan {
	opts := &readyOptions{
		pollInterval: DefaultPollInterval,
		timeout:      DefaultWaitTimeout,
		request:      nil,
	}
	for _, fn := range fns {
		fn(opts)
	}

	c := make(chan error)
	if opts.request == nil {
		close(c)
	}

	return c
}
