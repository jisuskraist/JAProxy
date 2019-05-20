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
	ipLimit       rate.Limit
	pathLimit     rate.Limit
	burst         int
	age           time.Duration
	sweepInterval time.Duration
	limiters      map[Type]map[string]*client
	l             sync.Mutex
}

func (m *MemLimiter) addLimiter(t Type, value string) *rate.Limiter {
	log.Info("adding client", value)
	var l *rate.Limiter

	switch t {
	case IpAddress:
		l = rate.NewLimiter(m.ipLimit, m.burst)
	case URL:
		l = rate.NewLimiter(m.pathLimit, m.burst)
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

func (m *MemLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(int(m.ipLimit)))

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		url := r.RequestURI
		ipl := m.getLimiter(IpAddress, ip)
		pathl := m.getLimiter(URL, url)

		if ipl.Allow() == false || pathl.Allow() == false {
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
		for t, l := range m.limiters {
			for val, cl := range l {
				if time.Now().Sub(cl.lastSeen) > m.age*time.Second {
					delete(m.limiters[t], val)
				}
			}
		}
		m.l.Unlock()
	}
}

func NewMemLimiter(ipLimit, pathLimit, burst int, age, sweepInterval time.Duration) *MemLimiter {
	l := &MemLimiter{
		rate.Limit(ipLimit),
		rate.Limit(pathLimit),
		burst,
		age,
		sweepInterval,
		make(map[Type]map[string]*client),
		sync.Mutex{},
	}

	return l
}
