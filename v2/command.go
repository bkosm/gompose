package gompose

import (
	"bytes"
	"os/exec"
)

type cmdOutput string

func run(cmd exec.Cmd) (cmdOutput, error) {
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	return cmdOutput(out.String()), err
}
