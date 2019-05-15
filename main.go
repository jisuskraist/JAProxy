package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

var netClient *http.Client
var netTransport *http.Transport

func init() {
	fmt.Println("init")

	netTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxConnsPerHost:     0,
		MaxIdleConns:        2000,
		MaxIdleConnsPerHost: 200,
		IdleConnTimeout:     10 * time.Second, //Avoid piling up of connections in case of bad sync
		TLSHandshakeTimeout: 30 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}

func main() {
	targetURL, err := url.Parse("http://httpbin.org")
	if err != nil {
		log.Fatal(err)
	}

	proxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		//Overw
		req.Host = targetURL.Host
		req.URL.Host = targetURL.Host
		req.URL.Scheme = targetURL.Scheme
		req.RequestURI = ""
		//This could be synced to avoid having large file descriptors size < 5000 if needed
		resp, err := netClient.Do(req)
		//If there was an error making the request, return 500 with err as body
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(rw, err)
			return
		}
		//Copying headers from response
		for key, values := range resp.Header {
			for _, value := range values {
				rw.Header().Set(key, value)
			}
		}
		rw.WriteHeader(resp.StatusCode)
		_, err = io.Copy(rw, resp.Body)

		//Close body to avoid leak and enable TCP connection reuse
		err = resp.Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	})

	err = http.ListenAndServe(":8000", proxy)

	if err != nil {
		panic(err)
	}
}
