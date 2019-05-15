package main

import (
	"github.com/jisuskraist/JAProxy/providers"
	"github.com/jisuskraist/JAProxy/services"
	"github.com/jisuskraist/JAProxy/structs"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func main() {
	log.SetLevel(log.DebugLevel)
	//create a channel to handle shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	log.Info("Starting app")

	prov := providers.NewJSONProvider("config.json")
	serv := structs.Services{}
	conf := structs.Config{}

	conf.LoadCommon(prov)
	conf.LoadNetwork(prov)
	conf.LoadRoutes(prov)

	http := &services.HTTPProxy{Config: conf}
	serv.AddService(http)

	serv.Run()

	//Keep running until a signal is received on stop channel.
	for {
		select {
		case <-stop:
			log.Info("Shutting down :(")
			serv.Stop() //Stop all running services
			return
		default:

		}
	}
}
