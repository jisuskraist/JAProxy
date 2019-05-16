package balancing

import (
	"github.com/jisuskraist/JAProxy/pkg/config"
	"net/url"
)

type Strategy int

const (
	RoundRobin = iota
)

type Balancer interface {
	NextTarget(host string) (*url.URL, error)
}

func NewBalanceStrategy(s Strategy, routes []config.RouteMapping) Balancer {
	switch s {
	case RoundRobin:
		return newRobinWood(routes)
	default:
		return newRobinWood(routes)
	}
}
