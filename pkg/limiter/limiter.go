package limiter

import (
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

type StorageType int

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
	addClient(ip string) *rate.Limiter
	getClient(ip string) *rate.Limiter
}

func NewLimiter(s StorageType, limit, burst int, age, sweepInterval time.Duration) Limiter {
	switch s {
	case InMemory:
		fallthrough
	default:
		return NewMemLimiter(limit, burst, age, sweepInterval)
	}
}
