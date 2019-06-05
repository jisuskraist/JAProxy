package main

import (
	"github.com/jisuskraist/JAProxy/pkg/balance"
	"github.com/jisuskraist/JAProxy/pkg/config"
	"github.com/jisuskraist/JAProxy/pkg/limiter"
	"github.com/jisuskraist/JAProxy/pkg/metrics"
	"github.com/jisuskraist/JAProxy/pkg/network"
	"github.com/jisuskraist/JAProxy/pkg/proxies"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Info("Starting app")

	// Configuration
	prov, err := config.NewProvider(config.Consul)
	if err != nil {
		panic(err)
	}
	conf := config.Config{}

	conf.LoadCommon(prov)
	conf.LoadNetwork(prov)
	conf.LoadRoutes(prov)

	log.Debugf("%+v\n", conf)
	// Metrics
	r := metrics.NewRegistry()
	_ = prometheus.Register(r.ProxyHist)
	// Limiter
	l := limiter.NewLimiter(limiter.Redis, conf.Limiter)
	go l.CleanUp()
	// Proxy
	proxy := proxies.NewHTTPProxy(conf.Client, balance.NewBalancer(balance.RoundRobin, conf.Routes), r)
	proxy.OnRequest(func(req *http.Request) {

	})
	proxy.OnResponse(func(res *http.Response) {

	})
	// HTTP server
	handler := http.NewServeMux()
	handler.Handle("/metrics", promhttp.Handler())
	handler.HandleFunc("/health", func(rw http.ResponseWriter, req *http.Request) {
		if l.IsHealthy() {
			rw.WriteHeader(http.StatusOK)
		} else {
			rw.WriteHeader(http.StatusServiceUnavailable)
		}

	})
	handler.Handle("/", l.Limit(
		network.RequestHandler{
			M: r.ProxyHist,
			H: proxy.ServeHTTP},
	))

	// Serving in HTTP, to serve HTTPS add certificates and spin up the server with them.
	err = http.ListenAndServe(":"+strconv.Itoa(conf.Port), handler)

	if err != nil {
		panic(err)
	}
}
