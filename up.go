package gompose

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type upOpts struct {
	wait           ReadyOrErrChan
	onSignal       func(os.Signal)
	customServices []string
	customFile     *string
}

type UpOption func(*upOpts)

func WithWait(c ReadyOrErrChan) UpOption {
	return func(o *upOpts) {
		o.wait = c
	}
}

func WithSignalCallback(fn func(os.Signal)) UpOption {
	return func(o *upOpts) {
		o.onSignal = fn
	}
}

func WithCustomServices(services ...string) UpOption {
	return func(o *upOpts) {
		o.customServices = services
	}
}

func AsUpOpt(fns ...GomposeOption) UpOption {
	g := &gomposeOpts{customFile: nil}
	for _, fn := range fns {
		fn(g)
	}

	return func(o *upOpts) {
		o.customFile = g.customFile
	}
}

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
