package network

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

type AppError struct {
	Error   error
	Message string
	Code    int
}

type RequestHandler struct {
	M *prometheus.HistogramVec
	H func(w http.ResponseWriter, r *http.Request) *AppError
}

func (rh RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	err := rh.H(w, r)
	duration := time.Since(start)
	if err != nil {
		rh.M.WithLabelValues(fmt.Sprintf("%d", err.Code)).Observe(duration.Seconds())
		http.Error(w, err.Message, err.Code)
		return
	}
	rh.M.WithLabelValues(fmt.Sprintf("%d", http.StatusOK)).Observe(duration.Seconds())
}
