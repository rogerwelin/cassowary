package client

import (
	"time"

	"go.uber.org/ratelimit"
)

type sleepLimited struct {
	d time.Duration
}

// NewSleepLimited returns a RateLimiter that is sleep limited.
func NewSleepLimited(d time.Duration) ratelimit.Limiter {
	return sleepLimited{d}
}

func (s sleepLimited) Take() time.Time {
	time.Sleep(s.d)

	return time.Now()
}
