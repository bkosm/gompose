package gompose

import (
	"bytes"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

type ReadyChan <-chan any

func ReadyOnStdout(cmd *exec.Cmd, awaiting string, times int) ReadyChan {
	c := make(chan any)

	go func() {
		for {
			if res, err := run(*cmd); err != nil {
				panic(res.error)
			} else {
				count := 0
				for _, line := range strings.Split(res.out, "\n") {
					if strings.Contains(line, awaiting) {
						count++
					}
				}

				if count >= times {
					close(c)
					return
				} else {
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()

	return c
}

func ReadyOnLog(text string, times int) ReadyChan {
	return ReadyOnStdout(
		exec.Command("docker-compose", "logs"),
		text,
		times,
	)
}

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
