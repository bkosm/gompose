package gompose

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// Up starts containers specified in compose file.
// It can be configured with a custom compose file path and a list of services to be started.
// Returns an error if shell command fails or if the provided channel returns an error.
// When provided WithWait option, the program execution is suspended until the channel is closed or returns an error.
func Up(opts ...Option) error {
	var (
		file customFile
		up   up
	)
	reduceUpOptions(&file, &up, opts)

	handleSignal(up.onSignal)

	args := getCommandArgs(string(file), up.customServices)
	if _, err := run(*exec.Command("docker-compose", args...)); err != nil {
		return err
	}

	return handleWait(up.wait)
}

func reduceUpOptions(file *customFile, u *up, opts []Option) {
	*u = up{
		wait:           nil,
		onSignal:       nil,
		customServices: nil,
	}

	for _, opt := range opts {
		if fn := opt.withCustomFileFunc; fn != nil {
			fn(file)
		}
		if fn := opt.withUpFunc; fn != nil {
			fn(u)
		}
	}
}

func getCommandArgs(customFile string, customServices []string) []string {
	var args []string
	if customFile != "" {
		args = []string{"-f", customFile}
	}
	args = append(args, "up", "-d")
	if customServices != nil {
		args = append(args, customServices...)
	}
	return args
}

func handleWait(c ReadyOrErrChan) error {
	if c != nil {
		return <-c
	}
	return nil
}

func handleSignal(callback func(os.Signal)) {
	if callback != nil {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			s := <-signalChan
			callback(s)
		}()
	}
}
