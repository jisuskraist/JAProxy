package limiter

import (
	"github.com/go-redis/redis"
	"github.com/go-redis/redis_rate"
	"github.com/jisuskraist/JAProxy/pkg/config"
	"net"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
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
		_, _, allowed := rl.limiter.Allow(ip, int64(rl.cfg.IpLimit), time.Second)
		if !allowed {
			h := w.Header()
			h.Set("X-RateLimit-Limit", strconv.Itoa(int(rl.cfg.IpLimit)))
			http.Error(w, "API rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CleanUp just to satisfy the interface, the library used automatically cleans the Keys already expired.
func (RedisLimiter) CleanUp() {

}

// IsHealthy pings the redis server to know if it's still responsive
func (rl RedisLimiter) IsHealthy() bool {
	res, err := rl.ring.Ping().Result()
	if err != nil {
		log.Error(res)
	}
	return err == nil
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
