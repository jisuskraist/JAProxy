package config

import (
	"crypto/tls"
	"encoding/json"
	"github.com/jisuskraist/JAProxy/pkg/balance"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

type JSONProvider struct {
	Network struct {
		Timeout               time.Duration `json:"timeout"`
		KeepAlive             time.Duration `json:"keepAlive"`
		MaxConnectionsPerHost int           `json:"maxConnectionsPerHost"`
		MaxIdleConns          int           `json:"maxIdleConns"`
		MaxIdleConnsPerHost   int           `json:"maxIdleConnsPerHost"`
		idleConnectionTimout  time.Duration `json:"idleConnectionTimeout"`
		TLSHandshakeTimeout   time.Duration `json:"TLSHandshakeTimeout"`
		TLSInsecureSkipVerify bool          `json:"TLSInsecureSkipVerify"`
	} `json:"network"`
	Common struct {
		ListenPort int `json:"listenPort"`
	} `json:"common"`
	Routes  []balance.RouteMapping `json:"routes"`
	Limiter LimiterConfig          `json:"limiter"`
}

//NewJSONProvider returns a new JSON provider for configuration.
func NewJSONProvider(path string) *JSONProvider {
	var config JSONProvider

	pwd, _ := os.Getwd()
	joinedPath := filepath.Join(pwd, path)
	log.Info("Loading config: " + joinedPath)
	configFile, err := os.Open(joinedPath)
	defer configFile.Close()

	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)

	if err != nil {
		panic(err)
	}

	log.Debug("Loaded configuration: ", config)
	return &config
}

func (p JSONProvider) LoadCommon(config *Config) {
	config.Port = p.Common.ListenPort
}

func (p JSONProvider) LoadRoutes(config *Config) {
	config.Routes = p.Routes
}

func (p JSONProvider) LoadNetwork(config *Config) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   p.Network.Timeout * time.Second,
			KeepAlive: p.Network.KeepAlive * time.Second,
		}).DialContext,
		MaxConnsPerHost:     p.Network.MaxConnectionsPerHost,
		MaxIdleConns:        p.Network.MaxIdleConns,
		MaxIdleConnsPerHost: p.Network.MaxIdleConnsPerHost,
		IdleConnTimeout:     p.Network.idleConnectionTimout * time.Second, //Avoid piling up of connections in case of bad sync
		TLSHandshakeTimeout: p.Network.TLSHandshakeTimeout * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: p.Network.TLSInsecureSkipVerify},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   p.Network.Timeout * time.Second,
	}

	config.Client = client
	config.Limiter = p.Limiter
}
