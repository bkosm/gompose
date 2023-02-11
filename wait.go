package gompose

import (
	"os/exec"
	"strings"
	"time"
)

type ReadyChan <-chan any

func ReadyOnStdout(cmd *exec.Cmd, awaiting string, times int) ReadyChan {
	c := make(chan any)

	go func() {
		for {
			if res, err := run(*cmd); err != nil {
				panic(err)
			} else {
				count := 0
				for _, line := range strings.Split(string(res), "\n") {
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
