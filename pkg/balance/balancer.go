package balance

import (
	"net/url"
)

type Strategy int

const (
	RoundRobin = iota
)

// Balancer defines NextTarget method which all balancers should call
// in order to get the next target in the list.
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
