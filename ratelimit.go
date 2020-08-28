package httpclient

import (
	"time"
)

type (
	RateLimiter interface {
		Take()
	}

	rateLimit struct {
		isTickedFirstTime bool
		duration          time.Duration
		ticker            *time.Ticker
	}

	nopRateLimit struct{}
)

func NewRateLimit(requests int64, period time.Duration) *rateLimit {
	return &rateLimit{
		duration: period / time.Duration(requests),
	}
}

func (r *rateLimit) Take() {
	if !r.isTickedFirstTime {
		r.isTickedFirstTime = true
		r.ticker = time.NewTicker(r.duration)
		return
	}

	<-r.ticker.C
}

func NewNopRateLimit() *nopRateLimit {
	return &nopRateLimit{}
}

func (r *nopRateLimit) Take() {
	return
}
