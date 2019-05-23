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
}

func NewLimiter(s StorageType, cfg config.LimiterConfig) Limiter {
	switch s {
	case Redis:
		servers := map[string]string{
			"server1": "localhost:6379",
		}
		return NewRedisLimiter(redisConfig{cfg, servers})
	case InMemory:
		fallthrough
	default:
		return NewMemLimiter(cfg)
	}
}
