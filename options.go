package httpclient

import (
	"net/http"
	"time"
)

type (
	Options struct {
		log         Logger
		timeout     time.Duration
		scheme      string
		rateLimit   RateLimiter
		loggerLevel LoggerLevel

		requestMetricHandler func(
			req *http.Request, header http.Header,
			h func(req *http.Request, header http.Header) (statusCode int, body []byte, err error),
		) (statusCode int, body []byte, err error)

		latencyMetricHandler func(
			req *http.Request, header http.Header,
			h func(req *http.Request, header http.Header) (statusCode int, body []byte, err error),
		) (statusCode int, body []byte, err error)
	}

	Option func(*Options)
)

func Log(log Logger) Option {
	return func(args *Options) {
		args.log = log
	}
}

func LogLevel(level LoggerLevel) Option {
	return func(args *Options) {
		args.loggerLevel = level
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

func RequestMetricHandler(
	metricHandler func(
		req *http.Request, header http.Header,
		h func(req *http.Request, header http.Header) (statusCode int, body []byte, err error),
	) (statusCode int, body []byte, err error),
) Option {
	return func(args *Options) {
		args.requestMetricHandler = metricHandler
	}
}

func LatencyMetricHandler(
	metricHandler func(
		req *http.Request, header http.Header,
		h func(req *http.Request, header http.Header) (statusCode int, body []byte, err error),
	) (statusCode int, body []byte, err error),
) Option {
	return func(args *Options) {
		args.latencyMetricHandler = metricHandler
	}
}
