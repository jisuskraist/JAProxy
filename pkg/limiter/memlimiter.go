package limiter

import (
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type MemLimiter struct {
	limit         rate.Limit
	burst         int
	age           time.Duration
	sweepInterval time.Duration
	clients       map[string]*client
	l             sync.Mutex
}

func (m *MemLimiter) addClient(ip string) *rate.Limiter {
	log.Info("adding client", ip)
	l := rate.NewLimiter(m.limit, m.burst)
	m.l.Lock()
	defer m.l.Unlock()

	m.clients[ip] = &client{l, time.Now()}
	return l
}

func (m *MemLimiter) getClient(ip string) *rate.Limiter {
	m.l.Lock()
	c, e := m.clients[ip]
	if !e {
		m.l.Unlock()
		return m.addClient(ip)
	}
	c.lastSeen = time.Now()
	m.l.Unlock()
	return c.limiter
}

func (m *MemLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(int(m.limit)))

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		l := m.getClient(ip)
		if l.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			log.Warn("Hit request limit")

			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *MemLimiter) CleanUp() {
	for {
		time.Sleep(m.sweepInterval * time.Second)
		m.l.Lock()
		for ip, v := range m.clients {
			if time.Now().Sub(v.lastSeen) > m.age*time.Second {
				delete(m.clients, ip)
			}
		}
		m.l.Unlock()
	}
}

func NewMemLimiter(limit, burst int, age, sweepInterval time.Duration) *MemLimiter {

	l := &MemLimiter{
		rate.Limit(limit),
		burst,
		age,
		sweepInterval,
		make(map[string]*client),
		sync.Mutex{},
	}

	return l
}
