package limiter

import (
	"github.com/jisuskraist/JAProxy/pkg/config"
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

// Limiter types
const (
	InMemory = iota
	Redis
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Interface defining the behaviour of a limiter
type Limiter interface {
	CleanUp()
	IsHealthy() bool
	Limit(next http.Handler) http.Handler
}

// NewLimiter returns a new limiter given the storage type
func NewLimiter(s StorageType, cfg config.LimiterConfig) Limiter {
	switch s {
	case Redis:
		return NewRedisLimiter(cfg)
	case InMemory:
		fallthrough
	default:
		return NewMemLimiter(cfg)
	}
}
