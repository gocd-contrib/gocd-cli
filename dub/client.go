package dub

import (
	"net"
	"net/http"
	"time"
)

type Client struct {
	native *http.Client
}

func (c *Client) Get(url string) *Request {
	return c.buildRequest("GET", url)
}

func (c *Client) Head(url string) *Request {
	return c.buildRequest("HEAD", url)
}

func (c *Client) Delete(url string) *Request {
	return c.buildRequest("DELETE", url)
}

func (c *Client) Put(url string) *Request {
	return c.buildRequest("PUT", url)
}

func (c *Client) Patch(url string) *Request {
	return c.buildRequest("PATCH", url)
}

func (c *Client) Post(url string) *Request {
	return c.buildRequest("POST", url)
}

func (c *Client) Connect(url string) *Request {
	return c.buildRequest("CONNECT", url)
}

func (c *Client) Trace(url string) *Request {
	return c.buildRequest("TRACE", url)
}

func (c *Client) Options(url string) *Request {
	return c.buildRequest("OPTIONS", url)
}

func (c *Client) buildRequest(method, url string) *Request {
	return &Request{Url: url, Method: method, c: c}
}

// Returns a dub.Client instance with a preconfigured http.Transport
func New() *Client {
	return Make(&http.Transport{
		Dial: (&net.Dialer{
			Timeout: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	})
}

// Returns a dub.Client instance allowing the user to specify their own http.Transport
func Make(t http.RoundTripper) *Client {
	return Wrap(&http.Client{Transport: t})
}

// Wraps an existing http.Client with a dub.Client
func Wrap(c *http.Client) *Client {
	return &Client{native: c}
}
