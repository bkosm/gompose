package gompose

import (
	"os"
	"os/exec"
	"os/signal"
)

func Up(ready ReadyChan, onInterrupt ...func()) error {
	if len(onInterrupt) > 0 {
		cleanup := onInterrupt[0]

		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, os.Interrupt)
		go func() {
			<-signalChan
			cleanup()
		}()
	}

	if _, err := run(*exec.Command("docker-compose", "up", "-d")); err != nil {
		return err
	}

	if ready != nil {
		_ = <-ready
	}
	return nil
}

func Down() error {
	_, err := run(*exec.Command("docker-compose", "down"))
	return err
}
