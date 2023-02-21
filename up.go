package gompose

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type (
	// UpOption  is a function that configures docker-compose up invocation.
	UpOption func(*upOpts)

	upOpts struct {
		wait           ReadyOrErrChan
		onSignal       func(os.Signal)
		customServices []string
		customFile     *string
	}
)

// WithWait suspends the programs' execution until the provided channel returns an error or is closed.
// See ReadyOrErrChan for more information.
func WithWait(c ReadyOrErrChan) UpOption {
	return func(o *upOpts) {
		o.wait = c
	}
}

// WithSignalCallback registers a callback function that is called when the program receives a SIGINT or SIGTERM signal.
// Useful for graceful shutdown of the program.
func WithSignalCallback(fn func(os.Signal)) UpOption {
	return func(o *upOpts) {
		o.onSignal = fn
	}
}

// WithCustomServices allows to declare a list of services specified in the compose file to be started.
func WithCustomServices(services ...string) UpOption {
	return func(o *upOpts) {
		o.customServices = services
	}
}

// AsUpOpt converts global GlobalOption which are useful in the context of Up to a UpOption.
func AsUpOpt(fns ...GlobalOption) UpOption {
	g := &globalOpts{customFile: nil}
	for _, fn := range fns {
		fn(g)
	}

	return func(o *upOpts) {
		o.customFile = g.customFile
	}
}

// Up starts containers specified in compose file.
// It can be configured with a custom compose file path and a list of services to be started.
// Returns an error if shell command fails or if the provided channel returns an error.
// When provided WithWait option, the program execution is suspended until the channel is closed or returns an error.
func Up(fns ...UpOption) error {
	opts := &upOpts{
		wait:           nil,
		onSignal:       nil,
		customServices: nil,
		customFile:     nil,
	}
	for _, fn := range fns {
		fn(opts)
	}

	handleSignal(opts.onSignal)

	args := getCommandArgs(opts.customFile, opts.customServices)
	if _, err := run(*exec.Command("docker-compose", args...)); err != nil {
		return err
	}

	return handleWait(opts.wait)
}

func getCommandArgs(customFile *string, customServices []string) []string {
	var args []string
	if customFile != nil {
		args = []string{"-f", *customFile}
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
