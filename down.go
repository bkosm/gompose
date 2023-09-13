package gompose

import (
	"os"
	"os/exec"
	"strings"
)

// Down stops and removes containers, networks, images, and volumes specified in compose file.
// It can be configured with a custom compose file path.
// Returns an error if shell command fails.
// Skips the command invocation if the current shell has a CSV environment variable with the key of SkipEnv
// set with a value of SkipDown. This allows for retaining the services between runs without altering source code.
func Down(opts ...Option) error {
	if shouldSkipCommand() {
		return nil
	}

	customFile := reduceCustomFileOptions(opts)

	var args []string
	if customFile != "" {
		args = []string{"-f", string(customFile)}
	}
	args = append(args, "down")

	_, err := run(*exec.Command("docker-compose", args...))
	return err
}

func shouldSkipCommand() bool {
	val, ok := os.LookupEnv(SkipEnv)
	if !ok {
		return false
	}

	for _, token := range strings.Split(val, ",") {
		if strings.ToLower(token) == SkipDown {
			return true
		}
	}
	return false
}
