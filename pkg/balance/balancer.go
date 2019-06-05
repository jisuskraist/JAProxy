package balance

import (
	"net/url"
)

type Strategy int

//RouteMapping represents a mapping of a domain with it's targets destinations.
//There could be one or more targets.
type RouteMapping struct {
	Domain  string
	Targets []string
}

// Balancer strategies
const (
	RoundRobin = iota
)

// Balancer defines NextTarget which all balancers should implement
// in order to get the next target in the list of targets.
type Balancer interface {
	NextTarget(host string) (*url.URL, error)
}

// NewBalancer returns a new balancer given a type and routes
func NewBalancer(s Strategy, routes []RouteMapping) Balancer {
	switch s {
	case RoundRobin:
		return newRobinWood(routes)
	default:
		return newRobinWood(routes)
	}
}
