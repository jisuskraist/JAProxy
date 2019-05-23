package limiter

import (
	"github.com/go-redis/redis"
	"github.com/go-redis/redis_rate"
	"github.com/jisuskraist/JAProxy/pkg/config"
	"net"
	"net/http"
	"time"
)

type RedisLimiter struct {
	ring    *redis.Ring
	limiter *redis_rate.Limiter
	cfg     redisConfig
}

type redisConfig struct {
	config.LimiterConfig
	Servers map[string]string
}

func (rl RedisLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		rate, _, allowed := rl.limiter.Allow(ip, int64(rl.cfg.IpLimit), time.Minute)
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

func (RedisLimiter) CleanUp() {

}

func NewRedisLimiter(cfg redisConfig) *RedisLimiter {
	r := redis.NewRing(&redis.RingOptions{
		Addrs: cfg.Servers,
	})
	
	return &RedisLimiter{
		ring:    r,
		limiter: redis_rate.NewLimiter(r),
	}
}
