package gompose

import (
	"log"
	"os/exec"
	"strings"
	"time"
)

type ReadyOnStdoutOption func(*readyOnStdoutOptions)

type readyOnStdoutOptions struct {
	awaiting     string
	times        uint
	timeout      time.Duration
	pollInterval time.Duration
}

func AwaitingText(text string) ReadyOnStdoutOption {
	return func(o *readyOnStdoutOptions) {
		o.awaiting = text
	}
}

func Times(n uint) ReadyOnStdoutOption {
	return func(o *readyOnStdoutOptions) {
		o.times = n
	}
}

func WithTimeout(t time.Duration) ReadyOnStdoutOption {
	return func(o *readyOnStdoutOptions) {
		o.timeout = t
	}
}

func WithPollInterval(t time.Duration) ReadyOnStdoutOption {
	return func(o *readyOnStdoutOptions) {
		o.pollInterval = t
	}
}

func ReadyOnStdout(cmd *exec.Cmd, fns ...ReadyOnStdoutOption) ReadyOrErrChan {
	opts := &readyOnStdoutOptions{
		awaiting:     "",
		times:        1,
		pollInterval: 100 * time.Millisecond,
		timeout:      10 * time.Minute,
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

func seekCondition(cmd *exec.Cmd, opts *readyOnStdoutOptions, found chan error) {
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
					log.Print(line)
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
