package timing

import (
	"time"
)

type FrameLimiter struct {
	interval time.Duration
	last     time.Time
}

func NewFrameLimiter(fps int) *FrameLimiter {
	if fps <= 0 {
		fps = 60
	}
	return &FrameLimiter{interval: time.Second / time.Duration(fps)}
}

func (f *FrameLimiter) Allow() bool {
	now := time.Now()
	if f.last.IsZero() || now.Sub(f.last) >= f.interval {
		f.last = now
		return true
	}
	return false
}
