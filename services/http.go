package services

import (
	"context"
	"fmt"
	"github.com/jisuskraist/JAProxy/structs"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

//HTTPProxy represents an http proxy service.
type HTTPProxy struct {
	Config structs.Config
	server *http.Server
}

//Start starts the web server/proxy.
func (prx *HTTPProxy) Start() {
	log.Info("Starting HTTP Proxy")
	//TODO: add a strategy to the target pick, implement a pool of services maybe?
	targetURL, err := url.Parse(prx.Config.Routes[0].Targets[0])
	if err != nil {
		log.Fatal(err)
	}

	//handles the request made to the running web server
	proxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		//Overwrite host, scheme and requestUri of proxied request
		req.Host = targetURL.Host
		req.URL.Host = targetURL.Host
		req.URL.Scheme = targetURL.Scheme
		req.RequestURI = ""
		//This could be synced to use file descriptors < 5000 if needed
		resp, err := prx.Config.Network.NetClient.Do(req)
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
		//Set status code and copy bodies
		rw.WriteHeader(resp.StatusCode)
		_, err = io.Copy(rw, resp.Body)

		//Close body to avoid leak and enable TCP connection reuse
		err = resp.Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	})
	prx.server = &http.Server{
		Addr:    ":" + strconv.Itoa(prx.Config.Port),
		Handler: proxy}
	//Start the server and send it to oblivion... or just in a goroutine.
	go func() {
		err = prx.server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
}

//Stop stops the proxy.
func (prx HTTPProxy) Stop() {
	log.Info("Stopping HTTProxy")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	prx.server.Shutdown(ctx)
}
