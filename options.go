package httpclient

type (
	Options struct {
		log Logger
	}

	Option func(*Options)
)

func Log(log Logger) Option {
	return func(args *Options) {
		args.log = log
	}
}
