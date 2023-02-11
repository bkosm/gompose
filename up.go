package gompose

import (
	"os"
	"os/exec"
	"os/signal"
)

type upOpts struct {
	wait           ReadyOrErrChan
	onSignal       func(os.Signal)
	customServices []string
	customFile     *string
}

type UpOption func(*upOpts)

func WaitFor(c ReadyOrErrChan) UpOption {
	return func(o *upOpts) {
		o.wait = c
	}
}

func OnSignal(fn func(os.Signal)) UpOption {
	return func(o *upOpts) {
		o.onSignal = fn
	}
}

func WithCustomServices(services ...string) UpOption {
	return func(o *upOpts) {
		o.customServices = services
	}
}

func AsUpOpt(fn GomposeOption) UpOption {
	g := &gomposeOpts{customFile: nil}
	fn(g)

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

	handleWait(opts.wait)
	return nil
}

func getCommandArgs(customFile *string, customServices []string) []string {
	var args []string
	if customFile != nil {
		args = append(args, "-f", *customFile)
	}
	args = append(args, "up", "-d")
	if customServices != nil {
		args = append(args, customServices...)
	}
	return args
}

func handleWait(c ReadyOrErrChan) {
	if c != nil {
		_ = <-c
	}
}

func handleSignal(callback func(os.Signal)) {
	if callback != nil {
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, os.Interrupt, os.Kill)
		go func() {
			s := <-signalChan
			callback(s)
		}()
	}
}
