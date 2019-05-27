package limiter

import (
	"github.com/go-redis/redis"
	"github.com/go-redis/redis_rate"
	"github.com/jisuskraist/JAProxy/pkg/config"
	"net"
	"net/http"
	"time"
)

// RedisLimiter represents a limiter storing the IP/Path in a redis database.
type RedisLimiter struct {
	ring    *redis.Ring
	limiter *redis_rate.Limiter
	cfg     config.LimiterConfig
}

// Limit is a middleware which stops the request if rate is exceeded or continues down the chain.
func (rl RedisLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement filtering by path here
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		rate, _, allowed := rl.limiter.Allow(ip, int64(rl.cfg.IpLimit), time.Second)
		if !allowed {
			h := w.Header()
			h.Set("X-RateLimit-Limit", string(rl.cfg.IpLimit))
			h.Set("X-RageLimit-Remaining", string(rl.cfg.IpLimit-rate))
			http.Error(w, "API rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CleanUp just to satisfy the interface, the library used automatically cleans the Keys already free.
func (RedisLimiter) CleanUp() {

}

// IsHealthy pings the redis server to know if it's still responsive
func (rl RedisLimiter) IsHealthy() bool {
	return rl.ring.Ping().Err() != nil
}

func NewRedisLimiter(cfg config.LimiterConfig) *RedisLimiter {
	r := redis.NewRing(&redis.RingOptions{
		Addrs: cfg.RedisAddress,
	})

	return &RedisLimiter{
		ring:    r,
		limiter: redis_rate.NewLimiter(r),
		cfg:     cfg,
	}
}
