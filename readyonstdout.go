package gompose

import (
	"os/exec"
	"strings"
	"time"
)

func ReadyOnStdout(cmd *exec.Cmd, fns ...ReadyOption) ReadyOrErrChan {
	opts := &readyOptions{
		awaiting:     "",
		times:        1,
		pollInterval: DefaultPollInterval,
		timeout:      DefaultWaitTimeout,
	}
	for _, fn := range fns {
		fn(opts)
	}
	readyOrErr := make(chan error)

	go seekOrTimeout(opts.timeout, readyOrErr, func(found chan error) {
		seekCondition(cmd, opts, found)
	})

	return readyOrErr
}

func seekCondition(cmd *exec.Cmd, opts *readyOptions, found chan error) {
	for {
		select {
		case <-found:
			return
		default:
			if res, err := run(*cmd); err != nil {
				found <- err
				close(found)
				return
			} else {
				count := 0
				for _, line := range strings.Split(string(res), "\n") {
					if strings.Contains(line, opts.awaiting) {
						count++
					}
				}

				if count >= int(opts.times) {
					close(found)
					return
				} else {
					time.Sleep(opts.pollInterval)
				}
			}
		}
	}
}
