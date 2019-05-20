package limiter

import (
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

type StorageType int
type Type int

const (
	IpAddress = iota
	URL
)

const (
	InMemory = iota
	Redis
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type Limiter interface {
	CleanUp()
	Limit(next http.Handler) http.Handler
	addLimiter(t Type, v string) *rate.Limiter
	getLimiter(t Type, v string) *rate.Limiter
}

func NewLimiter(s StorageType, ipLimit, pathLimit, burst int, age, sweepInterval time.Duration) Limiter {
	switch s {
	case InMemory:
		fallthrough
	default:
		return NewMemLimiter(ipLimit, pathLimit, burst, age, sweepInterval)
	}
}
