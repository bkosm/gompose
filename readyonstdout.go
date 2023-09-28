package gompose

import (
	"os/exec"
	"strings"
)

// ReadyOnStdout returns a ReadyOrErrChan that is ready when the specified string is found in the
// stdout or stderr produced by the specified exec.Cmd.
// The channel will be closed immediately if no options are specified.
// An ErrWaitTimedOut will be returned if the timeout is reached.
// Times is defaulted to 1, the timeout and interval to DefaultWaitTimeout, DefaultPollInterval.
// The command will be run once per poll interval.
func ReadyOnStdout(cmd *exec.Cmd, awaiting string, opts ...Option) ReadyOrErrChan {
	options := timeBased{
		times:        1,
		timeout:      DefaultWaitTimeout,
		pollInterval: DefaultPollInterval,
	}
	for _, opt := range opts {
		if fn := opt.withTimeBasedFunc; fn != nil {
			fn(&options)
		}
	}

	readyOrErr := make(chan error)

	go seekOrTimeout(options.timeout, options.pollInterval, readyOrErr, func() (bool, error) {
		if res, err := run(*cmd); err != nil {
			return false, err
		} else {
			return countLogOccurrences(res, awaiting) >= int(options.times), nil
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
