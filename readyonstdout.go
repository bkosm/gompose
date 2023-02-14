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
	c := make(chan error)

	go func() {
		found := make(chan error)
		go seekCondition(cmd, opts, found)

		select {
		case err := <-found:
			if err != nil {
				c <- err
			}
			close(c)
			return
		case <-time.After(opts.timeout):
			c <- ErrWaitTimedOut
			close(c)
			close(found)
			return
		}
	}()

	return c
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
