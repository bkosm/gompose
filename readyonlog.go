package gompose

import "os/exec"

// ReadyOnLog returns a ReadyOrErrChan that is ready when the specified log message is found in the
// aggregated logs of all composed services.
// The channel will be closed immediately if no options are specified.
// An ErrWaitTimedOut will be returned if the timeout is reached.
// Times is defaulted to 1.
func ReadyOnLog(fns ...ReadyOption) ReadyOrErrChan {
	opts := &readyOptions{customFile: nil}
	for _, fn := range fns {
		fn(opts)
	}

	var args []string
	if opts.customFile != nil {
		args = []string{"-f", *opts.customFile}
	}
	args = append(args, "logs")

	return ReadyOnStdout(
		exec.Command("docker-compose", args...),
		fns...,
	)
}
