package gompose

import (
	"os/exec"
	"strings"
)

// ReadyOnStdout returns a ReadyOrErrChan that is ready when the specified string is found in the
// stdout or stderr produced by the specified exec.Cmd.
// The channel will be closed immediately if no options are specified.
// An ErrWaitTimedOut will be returned if the timeout is reached.
// Times is defaulted to 1.
// The command will be run once per poll interval.
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
