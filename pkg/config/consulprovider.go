package config

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/consul/api"
	"os"
)

// ConsulProvider is a provider which populates the configuration from a Consul Server.
type ConsulProvider struct {
	json   *JSONProvider
	client api.Client
}

func (cp ConsulProvider) LoadNetwork(config *Config) {
	cp.json.LoadNetwork(config)
}

func (cp ConsulProvider) LoadCommon(config *Config) {
	cp.json.LoadCommon(config)
}

func (cp ConsulProvider) LoadRoutes(config *Config) {
	cp.json.LoadRoutes(config)
}

func (cp ConsulProvider) LoadAll(config *Config) {
	cp.json.LoadAll(config)
}

// NewConsulProvider returns a new consul configuration provider.
// Since we store the configuration as a JSON in Consul to cut off dev times
// we end up returning a JSON provider embedded, but this should be refactored accordingly.
func NewConsulProvider() (cfg *JSONProvider, err error) {
	c, err := api.NewClient(&api.Config{
		Address: os.Getenv("CONSUL_ADDR"),
	})
	pair, _, err := c.KV().Get("PROXY_CONFIG", nil)
	if err != nil {
		return
	}
	if pair == nil {
		return nil, errors.New("invalid key/pair")
	}
	cp := &ConsulProvider{}
	err = json.Unmarshal(pair.Value, &cp.json)
	return cp.json, err
}
