package gompose

import (
	"os/exec"
)

type DownOption func(*downOpts)

type downOpts struct {
	customFile *string
}

func AsDownOpt(fns ...GomposeOption) DownOption {
	g := &gomposeOpts{customFile: nil}
	for _, fn := range fns {
		fn(g)
	}

	return func(o *downOpts) {
		o.customFile = g.customFile
	}
}

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
