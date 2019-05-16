package services

import (
	"fmt"
	"github.com/jisuskraist/JAProxy/pkg/config"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

type reqHandler func(req *http.Request)
type resHandler func(res *http.Response)

//HTTPProxy represents an http proxy service.
type HTTPProxy struct {
	cfg         config.Config
	reqHandlers []reqHandler
	resHandler  []resHandler
}

// Handles the request made to the web server.
func (p *HTTPProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	p.requestMiddleware(req)

	//TODO: add a strategy to the target pick, implement a pool of services maybe?
	targetURL, err := url.Parse(p.cfg.Routes[0].Targets[0])
	if err != nil {
		log.Fatal(err)
	}

	//Overwrite host, scheme and requestUri of proxied request
	overwriteRequest(req, *targetURL)
	//This could be synced to use less file descriptors if needed
	resp, err := p.cfg.Network.NetClient.Do(req)
	//If there was an error making the request, return 500 with err as body
	if err != nil {
		log.Error("An error occurred while sending request " , err)
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

func NewHTTPProxy(cfg config.Config) *HTTPProxy {
	proxy := &HTTPProxy{
		cfg: cfg,
	}

	return proxy
}
