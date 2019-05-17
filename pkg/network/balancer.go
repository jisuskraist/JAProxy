package network

import (
	"net/url"
)

type Strategy int

const (
	RoundRobin = iota
)

type Balancer interface {
	NextTarget(host string) (*url.URL, error)
}

func NewBalancer(s Strategy, routes []RouteMapping) Balancer {
	switch s {
	case RoundRobin:
		return newRobinWood(routes)
	default:
		return newRobinWood(routes)
	}
}
