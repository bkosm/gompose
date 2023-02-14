package gompose

import (
	"os/exec"
	"strings"
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

	go seekOrTimeout(opts.timeout, opts.pollInterval, readyOrErr, func() (bool, error) {
		if res, err := run(*cmd); err != nil {
			return false, err
		} else {
			return countLogOccurrences(res, opts.awaiting) >= int(opts.times), nil
		}
	})

	return readyOrErr
}

func countLogOccurrences(res cmdOutput, awaiting string) int {
	count := 0
	for _, line := range strings.Split(string(res), "\n") {
		if strings.Contains(line, awaiting) {
			count++
		}
	}
	return count
}
