package limiter

import (
	"github.com/jisuskraist/JAProxy/pkg/config"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// MemLimiter implements a simple limiter storing the visitors in memory
type MemLimiter struct {
	cfg           config.LimiterConfig
	limiters      map[Type]map[string]*client
	l             sync.Mutex
}

func (m *MemLimiter) addLimiter(t Type, value string) *rate.Limiter {
	log.Info("adding client", value)
	var l *rate.Limiter

	switch t {
	case IpAddress:
		l = rate.NewLimiter(rate.Limit(m.cfg.IpLimit), m.cfg.Burst)
	case URL:
		l = rate.NewLimiter(rate.Limit(m.cfg.PathLimit), m.cfg.Burst)
	}

	m.l.Lock()
	defer m.l.Unlock()
	cls := make(map[string]*client)
	cls[value] = &client{l, time.Now()}
	m.limiters[t] = cls
	return l
}

func (m *MemLimiter) getLimiter(t Type, val string) *rate.Limiter {
	m.l.Lock()
	c, e := m.limiters[t][val]
	if !e {
		m.l.Unlock()
		return m.addLimiter(t, val)
	}
	c.lastSeen = time.Now()
	m.l.Unlock()
	return c.limiter
}

// Limit is a middleware which stops the request if rate is exceeded or continues down the chain.
func (m *MemLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(int(m.cfg.IpLimit)))

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		url := r.RequestURI
		ipl := m.getLimiter(IpAddress, ip)
		pathl := m.getLimiter(URL, url)

		if ipl.Allow() == false || pathl.Allow() == false {
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(int(m.cfg.IpLimit)))
			http.Error(w, "API rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CleanUp cleans the memory to avoid insane memory grow
func (m *MemLimiter) CleanUp() {
	for {
		time.Sleep(m.cfg.SweepInterval * time.Second)
		m.l.Lock()
		for t, l := range m.limiters {
			for val, cl := range l {
				if time.Now().Sub(cl.lastSeen) > m.cfg.Age*time.Second {
					delete(m.limiters[t], val)
				}
			}
		}
		m.l.Unlock()
	}
}

// IsHealthy returns if the current memory limiter is healthy
// Unless we crashed memory should be alright, right? Guys...?
func (MemLimiter) IsHealthy() bool {
	return true
}

func NewMemLimiter(cfg config.LimiterConfig) *MemLimiter {
	l := &MemLimiter{
		cfg,
		make(map[Type]map[string]*client),
		sync.Mutex{},
	}

	return l
}
