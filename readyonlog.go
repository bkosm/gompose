package gompose

import "os/exec"

func ReadyOnLog(fns ...ReadyOption) ReadyOrErrChan {
	opts := &readyOptions{
		customFile: nil,
	}
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
