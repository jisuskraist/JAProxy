package structs

import "net/http"

//ConfigurationProvider defines an interface to be
//implemented by all the configurations providers.
//This enables to have many configuration providers such as files, consul, etc.
type ConfigurationProvider interface {
	LoadCommon(config *Config)
	LoadRoutes(config *Config)
	LoadNetwork(config *Config)
}

//RouteMapping represents a mapping of a domain with it's targets destinations.
//There could be one or more targets.
type RouteMapping struct {
	Domain  string
	Targets []string
}

//Network represents the network configuration and the needed structs.
type Network struct {
	NetClient *http.Client
}

//Config hold the configuration of an application such as routes, listen port, network configuration.
type Config struct {
	Routes  []RouteMapping
	Port    int
	Network Network
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
