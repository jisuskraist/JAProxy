package proxies

import (
	"github.com/jisuskraist/JAProxy/pkg/balance"
	"github.com/jisuskraist/JAProxy/pkg/metrics"
	"github.com/jisuskraist/JAProxy/pkg/network"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type reqHandler func(req *http.Request)
type resHandler func(res *http.Response)

//HTTPProxy represents an http proxy service.
type HTTPProxy struct {
	balancer    balance.Balancer
	netClient   network.Client
	registry    *metrics.Registry
	reqHandlers []reqHandler
	resHandler  []resHandler
}

// Handles the request made to the web server.
func (p *HTTPProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) *network.AppError {
	p.requestMiddleware(req)

	targetURL, err := p.balancer.NextTarget(req.Host)
	//TODO: this error handling should be delegated to an error handling middleware to avoid duped code
	if err != nil {
		log.Warn(err)
		return &network.AppError{
			Message: "Could not find a target host",
			Code:    http.StatusBadGateway,
			Error:   err,
		}
	}
	h, _, _ := net.SplitHostPort(req.RemoteAddr)
	req.Header.Set("X-Forwarded-For", h)
	//Overwrite host, scheme and requestUri of proxied request
	overwriteRequest(req, *targetURL)
	//This could be synced to use less file descriptors if needed
	resp, err := p.netClient.Do(req)
	//If there was an error making the request, return 500 with err as body
	if err != nil {
		log.Error("An error occurred while sending request ", err)
		return &network.AppError{
			Message: "Could not reach target",
			Code:    http.StatusInternalServerError,
			Error:   err,
		}
	}
	p.responseMiddleware(resp)
	copyHeaders(rw.Header(), resp.Header, true)
	copyCookies(rw.Header(), resp.Cookies())
	//Set status code and copy bodies
	rw.WriteHeader(resp.StatusCode)
	if resp.Body != nil {
		//To support streams we flush a lot in a separate routine.
		done := make(chan bool)
		go func() {
			for {
				select {
				case <-time.Tick(50 * time.Millisecond):
					rw.(http.Flusher).Flush()
				case <-done:
					return
				}
			}
		}()
		//TODO: add support for trailers
		//https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Trailer
		_, err = io.Copy(rw, resp.Body)
		close(done) // Avoid leak
		if err != nil {
			return &network.AppError{
				Message: "An error occurred while copying bodies",
				Code:    http.StatusInternalServerError,
				Error:   err,
			}
		}
	}
	//Close body to avoid leak and enable TCP connection reuse (in case of using golang net client)
	err = resp.Body.Close()
	if err != nil {
		log.Error("An error occurred while closing the response body %s", err.Error())
	}
	return nil
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

func copyCookies(dst http.Header, src []*http.Cookie) {
	for _, c := range src {
		dst.Add("Set-Cookie", c.Raw)
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
// some manipulation can broke the response further down the stream.
func (p *HTTPProxy) OnResponse(fn func(r *http.Response)) {
	p.resHandler = append(p.resHandler, fn)
}

func NewHTTPProxy(n network.Client, b balance.Balancer, r *metrics.Registry) *HTTPProxy {
	proxy := &HTTPProxy{
		netClient: n,
		balancer:  b,
		registry:  r,
	}

	return proxy
}
