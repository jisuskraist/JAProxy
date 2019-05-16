package balancing

import (
	"errors"
	"github.com/jisuskraist/JAProxy/pkg/config"
	"net/url"
	"sync"
)

type RobinWood struct {
	counter int
	l       *sync.Mutex
	routes  []config.RouteMapping
}

func (rw *RobinWood) NextTarget(host string) (*url.URL, error) {
	rw.l.Lock()
	defer rw.l.Unlock()

	for _, r := range rw.routes {
		if r.Domain == host {
			rw.counter = (rw.counter + 1) % len(r.Targets)
			return url.Parse(r.Targets[rw.counter])
		}
	}
	return nil, errors.New("couldn't find a host to map")
}

func newRobinWood(routes []config.RouteMapping) *RobinWood {
	return &RobinWood{
		routes:  routes,
		counter: 0,
		l:       &sync.Mutex{},
	}
}
