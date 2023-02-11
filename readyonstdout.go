package gompose

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type ReadyChan <-chan any

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

func ReadyOnStdout(cmd *exec.Cmd, fns ...ReadyOnStdoutOption) ReadyChan {
	opts := &readyOnStdoutOptions{
		awaiting:     "",
		times:        0,
		pollInterval: 100 * time.Millisecond,
		timeout:      10 * time.Minute,
	}
	for _, fn := range fns {
		fn(opts)
	}
	c := make(chan any)

	go func() {
		found := make(chan any)
		go seekCondition(cmd, opts, found)

		select {
		case <-found:
			close(c)
			return
		case <-time.After(opts.timeout):
			panic(fmt.Sprintf("gompose: wait condition timed out after %v", opts.timeout))
		}
	}()

	return c
}

func seekCondition(cmd *exec.Cmd, opts *readyOnStdoutOptions, found chan any) {
	for {
		if res, err := run(*cmd); err != nil {
			panic(err)
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
