package network

import "net/http"

// Is a client capable of handling a request and returning a response
// It could be a http client, tcp client or even a custom protocol client.
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}
