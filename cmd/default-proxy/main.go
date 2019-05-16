package main

import (
	"github.com/jisuskraist/JAProxy/pkg/config"
	"github.com/jisuskraist/JAProxy/pkg/services"
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

	proxy := services.NewHTTPProxy(conf)

	proxy.OnRequest(func(req *http.Request) {
		log.Debug(req.URL)
	})

	proxy.OnResponse(func(res *http.Response) {

	})

	http.ListenAndServe(":"+strconv.Itoa(conf.Port), proxy)
}
