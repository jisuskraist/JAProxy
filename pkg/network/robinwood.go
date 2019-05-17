package network

import (
	"errors"
	"net/url"
	"sync"
)
// RobinWood is a simple round robin
type RobinWood struct {
	counter int
	l       *sync.Mutex
	routes  []RouteMapping
}

func (rw *RobinWood) NextTarget(host string) (*url.URL, error) {
	rw.l.Lock()
	defer rw.l.Unlock()

	//TODO: make target iteration domain specific, currently counter shared by domains
	for _, r := range rw.routes {
		if r.Domain == host {
			rw.counter = (rw.counter + 1) % len(r.Targets)
			return url.Parse(r.Targets[rw.counter])
		}
	}
	return nil, errors.New("couldn't find a host to map")
}

func newRobinWood(routes []RouteMapping) *RobinWood {
	return &RobinWood{
		routes:  routes,
		counter: 0,
		l:       &sync.Mutex{},
	}
}
