package gompose

type (
	// GlobalOption is a function that configures global gompose options, shared by all gompose commands.
	GlobalOption func(*globalOpts)

	globalOpts struct {
		customFile *string
	}
)

// WithCustomFile sets the path of a custom compose file to be used by gompose.
func WithCustomFile(filepath string) GlobalOption {
	return func(o *globalOpts) {
		o.customFile = &filepath
	}
}
