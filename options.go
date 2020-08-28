package httpclient

import "time"

type (
	Options struct {
		log       Logger
		timeout   time.Duration
		scheme    string
		rateLimit RateLimiter
	}

	Option func(*Options)
)

func Log(log Logger) Option {
	return func(args *Options) {
		args.log = log
	}
}

func Timeout(timeout time.Duration) Option {
	return func(args *Options) {
		args.timeout = timeout
	}
}

func Scheme(scheme string) Option {
	return func(args *Options) {
		args.scheme = scheme
	}
}

func RateLimit(requests int64, period time.Duration) Option {
	return func(args *Options) {
		args.rateLimit = NewRateLimit(requests, period)
	}
}
