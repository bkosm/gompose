package gompose

import (
	"bytes"
	"os/exec"
)

type cmdResult struct {
	error error
	out   string
}

func run(cmd exec.Cmd) (cmdResult, error) {
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	return cmdResult{error: err, out: out.String()}, err
}
