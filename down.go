package gompose

import (
	"os/exec"
)

// Down stops and removes containers, networks, images, and volumes specified in compose file.
// It can be configured with a custom compose file path.
// Returns an error if shell command fails.
func Down(opts ...Option) error {
	var customFile customFile
	reduceCustomFileOptions(&customFile, opts)

	var args []string
	if customFile != "" {
		args = []string{"-f", string(customFile)}
	}
	args = append(args, "down")

	_, err := run(*exec.Command("docker-compose", args...))
	return err
}
