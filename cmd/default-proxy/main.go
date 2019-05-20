package main

import (
	"github.com/jisuskraist/JAProxy/pkg/balance"
	"github.com/jisuskraist/JAProxy/pkg/config"
	"github.com/jisuskraist/JAProxy/pkg/limiter"
	"github.com/jisuskraist/JAProxy/pkg/metrics"
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
	prov, err := config.NewProvider(config.JSON)
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
	prometheus.Register(r.Histogram)
	// Limiter
	l := limiter.NewLimiter(limiter.InMemory, conf.Limiter.IpLimit, conf.Limiter.PathLimit, conf.Limiter.Burst, conf.Limiter.Age, conf.Limiter.SweepInterval)
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
	handler.HandleFunc("/", proxy.ServeHTTP)

	http.ListenAndServe(":"+strconv.Itoa(conf.Port), l.Limit(handler))
}
