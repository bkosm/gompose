package gompose

import "os/exec"

func ReadyOnLog(fns ...ReadyOnStdoutOption) ReadyOrErrChan {
	return ReadyOnStdout(
		exec.Command("docker-compose", "logs"),
		fns...,
	)
}
