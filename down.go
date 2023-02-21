package gompose

import (
	"os/exec"
)

type (
	// DownOption is a function that configures docker-compose down invocation.
	DownOption func(*downOpts)

	downOpts struct {
		customFile *string
	}
)

// AsDownOpt converts GlobalOption s which are useful in the context of Down to a DownOption.
func AsDownOpt(fns ...GlobalOption) DownOption {
	g := &globalOpts{customFile: nil}
	for _, fn := range fns {
		fn(g)
	}

	return func(o *downOpts) {
		o.customFile = g.customFile
	}
}

// Down stops and removes containers, networks, images, and volumes specified in compose file.
// It can be configured with a custom compose file path.
// Returns an error if shell command fails.
func Down(fns ...DownOption) error {
	opts := &downOpts{customFile: nil}
	for _, fn := range fns {
		fn(opts)
	}

	var args []string
	if opts.customFile != nil {
		args = []string{"-f", *opts.customFile}
	}
	args = append(args, "down")

	_, err := run(*exec.Command("docker-compose", args...))
	return err
}
