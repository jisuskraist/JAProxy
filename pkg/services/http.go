package services

import (
	"fmt"
	"github.com/jisuskraist/JAProxy/pkg/balancing"
	"github.com/jisuskraist/JAProxy/pkg/metrics"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

type reqHandler func(req *http.Request)
type resHandler func(res *http.Response)

//HTTPProxy represents an http proxy service.
type HTTPProxy struct {
	balancer    balancing.Balancer
	netClient   *http.Client
	reqHandlers []reqHandler
	resHandler  []resHandler
}

// Handles the request made to the web server.
func (p *HTTPProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	metrics.RpsCounter.Incr(1)

	p.requestMiddleware(req)

	targetURL, err := p.balancer.NextTarget(req.Host)

	if err != nil {
		log.Warn(err)
		rw.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(rw, err)
		return
	}

	//Overwrite host, scheme and requestUri of proxied request
	overwriteRequest(req, *targetURL)
	//This could be synced to use less file descriptors if needed
	resp, err := p.netClient.Do(req)
	//If there was an error making the request, return 500 with err as body
	if err != nil {
		log.Error("An error occurred while sending request ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(rw, err)
		return
	}
	p.responseMiddleware(resp)
	//Copying headers from response
	copyHeaders(rw.Header(), resp.Header, false)
	//Set status code and copy bodies
	rw.WriteHeader(resp.StatusCode)
	_, err = io.Copy(rw, resp.Body)

	if err != nil {
		log.Error("An error occurred while copying response from server %s", err.Error())
	}

	//Close body to avoid leak and enable TCP connection reuse
	err = resp.Body.Close()
	if err != nil {
		log.Error("An error occurred while closing the response body %s", err.Error())
	}
}

// Overwrites the request data to acts on its behalf prior to send it to the target server
func overwriteRequest(req *http.Request, target url.URL) {
	req.Host = target.Host
	req.URL.Host = target.Host
	req.URL.Scheme = target.Scheme
	req.RequestURI = ""
}

// Copy all headers from the target server response. The destination headers can be erase by setting the bool
func copyHeaders(dst, src http.Header, keepDestHeaders bool) {
	if !keepDestHeaders {
		for key := range dst {
			dst.Del(key)
		}
	}

	for key, values := range src {
		for _, value := range values {
			dst.Set(key, value)
		}
	}
}

func (p HTTPProxy) requestMiddleware(req *http.Request) (r *http.Request) {
	r = req
	for _, h := range p.reqHandlers {
		h(req)
	}
	return
}

func (p HTTPProxy) responseMiddleware(rw *http.Response) (r *http.Response) {
	r = rw
	for _, h := range p.resHandler {
		h(rw)
	}
	return
}

// Adds a callback to the request received by the proxy
func (p *HTTPProxy) OnRequest(fn func(r *http.Request)) {
	p.reqHandlers = append(p.reqHandlers, fn)
}

// Adds a callback for the response object. BEWARE! since this is a pointer
// some manipulation can broke the response further down the chain.
func (p *HTTPProxy) OnResponse(fn func(r *http.Response)) {
	p.resHandler = append(p.resHandler, fn)
}

func NewHTTPProxy(n *http.Client, b balancing.Balancer) *HTTPProxy {
	proxy := &HTTPProxy{
		netClient: n,
		balancer:  b,
	}

	return proxy
}
