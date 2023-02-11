package gompose

type GomposeOption func(*gomposeOpts)

type gomposeOpts struct {
	customFile *string
}

func WithCustomFile(filepath string) GomposeOption {
	return func(o *gomposeOpts) {
		o.customFile = &filepath
	}
}
