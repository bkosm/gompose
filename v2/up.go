package gompose

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// Up starts containers specified in compose file.
// Returns an error if shell command fails or if the provided channel returns an error.
// When provided CustomFile option, compose definitions from that file are used.
// When provided CustomServices option, only the specified services will be run from the compose spec.
// When provided Wait option, the program execution is suspended until the channel is closed or returns an error.
// When provided SignalCallback option, the specified function will be run on system interrupt.
// When provided RetryCommand option, the docker-compose command will be retried specified number of times in case of failure,
// each after the specified interval. Defaults to a single execution.
func Up(opts ...Option) error {
	var customFile customFile
	options := up{
		wait:           nil,
		onSignal:       nil,
		customServices: nil,
	}
	retry := retry{
		times: 1,
	}

	for _, opt := range opts {
		if fn := opt.withCustomFileFunc; fn != nil {
			fn(&customFile)
		}
		if fn := opt.withUpFunc; fn != nil {
			fn(&options)
		}
		if fn := opt.withRetryFunc; fn != nil {
			fn(&retry)
		}
	}

	handleSignal(options.onSignal)

	args := getCommandArgs(string(customFile), options.customServices)
	err := doRetry(retry.times, retry.interval, func() error {
		_, err := run(*exec.Command("docker-compose", args...))
		return err
	})
	if err != nil {
		return err
	}

	return handleWait(options.wait)
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
