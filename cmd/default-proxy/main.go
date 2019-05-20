package main

import (
	"github.com/jisuskraist/JAProxy/pkg/balance"
	"github.com/jisuskraist/JAProxy/pkg/config"
	"github.com/jisuskraist/JAProxy/pkg/limiter"
	"github.com/jisuskraist/JAProxy/pkg/proxies"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func main() {
	//log.SetLevel(log.DebugLevel)
	log.Info("Starting app")

	prov, err := config.NewProvider(config.JSON)
	if err != nil {
		panic(err)
	}
	conf := config.Config{}

	conf.LoadCommon(prov)
	conf.LoadNetwork(prov)
	conf.LoadRoutes(prov)

	l := limiter.NewLimiter(limiter.InMemory, 3, 5, 60, 180)
	go l.CleanUp()

	proxy := proxies.NewHTTPProxy(conf.Client, balance.NewBalancer(balance.RoundRobin, conf.Routes))

	proxy.OnRequest(func(req *http.Request) {
		log.Debug(req.URL)
	})

	proxy.OnResponse(func(res *http.Response) {

	})
	handler := http.NewServeMux()
	handler.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hey"))
	})
	handler.HandleFunc("/", proxy.ServeHTTP)

	http.ListenAndServe(":"+strconv.Itoa(conf.Port), l.Limit(handler))
}
