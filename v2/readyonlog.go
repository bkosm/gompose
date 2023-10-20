package gompose

import "os/exec"

// ReadyOnLog returns a ReadyOrErrChan that is ready when the specified log message is found in the
// aggregated logs of all composed services.
// This function proxies the responsibility to ReadyOnStdout by providing it the output of compose logs.
// Can be configured with CustomFile to read from custom compose spec.
func ReadyOnLog(awaiting string, opts ...Option) ReadyOrErrChan {
	var customFile customFile
	for _, opt := range opts {
		if fn := opt.withCustomFileFunc; fn != nil {
			fn(&customFile)
		}
	}

	var args []string
	if customFile != "" {
		args = []string{"-f", string(customFile)}
	}
	args = append(args, "logs")

	return ReadyOnStdout(
		exec.Command("docker-compose", args...),
		awaiting,
		opts...,
	)
}
