package metrics

import (
	"github.com/paulbellamy/ratecounter"
	"time"

	log "github.com/sirupsen/logrus"
)

var RpsCounter *ratecounter.RateCounter
//TODO: init should ve avoided ;)
func init() {
	RpsCounter = ratecounter.NewRateCounter(1 * time.Second)

	go func() {
		for {
			log.Info(RpsCounter.Rate())
			time.Sleep(1 * time.Second)
		}
	}()
}
