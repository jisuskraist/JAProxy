package config

import (
	"errors"
	"github.com/jisuskraist/JAProxy/pkg/network"
)

const JsonFile = "config.json"

type Type int

const (
	JSON Type = iota
)

//ConfigurationProvider defines an interface to be
//implemented by all the configurations providers.
//This enables to have many configuration providers such as
//files or consul and make them load specific parts of the configuration
type ConfigurationProvider interface {
	LoadCommon(config *Config)
	LoadRoutes(config *Config)
	LoadNetwork(config *Config)
}

func NewProvider(t Type) (ConfigurationProvider, error) {
	switch t {
	case JSON:
		return NewJSONProvider(JsonFile), nil
	default:
		return nil, errors.New("provider not defined")
	}
}

//Config hold the configuration of an application such as routes, listen port, network configuration.
type Config struct {
	Port   int
	Routes []network.RouteMapping
	Client network.Client
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
