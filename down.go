package gompose

import (
	"os/exec"
)

func Down() error {
	_, err := run(*exec.Command("docker-compose", "down"))
	return err
}
