package config

import (
	"errors"
	"github.com/jisuskraist/JAProxy/pkg/balance"
	"github.com/jisuskraist/JAProxy/pkg/network"
	"time"
)

const JsonFile = "/app/config/config.json"

type Type int

const (
	JSON Type = iota
	Consul
)

//ConfigurationProvider defines an interface to be
//implemented by all the configurations providers.
//This enables to have many configuration providers such as
//files or consul and make them load specific parts of the configuration
type ConfigurationProvider interface {
	LoadCommon(config *Config)
	LoadRoutes(config *Config)
	LoadNetwork(config *Config)
	LoadAll(config *Config)
}

func NewProvider(t Type) (ConfigurationProvider, error) {
	switch t {
	case JSON:
		return NewJSONProvider(JsonFile), nil
	case Consul:
		return NewConsulProvider()
	default:
		return nil, errors.New("provider not defined")
	}
}

//Config hold the configuration of an application such as routes, listen port, network configuration.
type Config struct {
	Port    int
	Routes  []balance.RouteMapping
	Client  network.Client
	Limiter LimiterConfig
}

// LimiterConfig is the configuration to be used by all the limiters.
type LimiterConfig struct {
	IpLimit       int64             `json:"ipLimit"`
	PathLimit     int64             `json:"pathLimit"`
	Burst         int               `json:"burst"`
	Age           time.Duration     `json:"age"`
	SweepInterval time.Duration     `json:"sweepInterval"`
	RedisAddress  map[string]string `json:"redisAddress"`
}

func (c *Config) LoadCommon(provider ConfigurationProvider) {
	provider.LoadCommon(c)
}

func (c *Config) LoadRoutes(provider ConfigurationProvider) {
	provider.LoadRoutes(c)
}

func (c *Config) LoadNetwork(provider ConfigurationProvider) {
	provider.LoadNetwork(c)
}

func (c *Config) LoadAll(provider ConfigurationProvider) {
	provider.LoadAll(c)
}
