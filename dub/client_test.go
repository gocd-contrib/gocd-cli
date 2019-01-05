package dub

import (
	"testing"
)

func TestClientRequestMethods(t *testing.T) {
	c := nopCl()
	as := asserts(t)

	get := c.Get("http://foo.bar")
	as.eq("GET", get.Method)
	as.eq("http://foo.bar", get.Url)

	post := c.Post("http://foo.bar")
	as.eq("POST", post.Method)
	as.eq("http://foo.bar", post.Url)

	put := c.Put("http://foo.bar")
	as.eq("PUT", put.Method)
	as.eq("http://foo.bar", put.Url)

	patch := c.Patch("http://foo.bar")
	as.eq("PATCH", patch.Method)
	as.eq("http://foo.bar", patch.Url)

	head := c.Head("http://foo.bar")
	as.eq("HEAD", head.Method)
	as.eq("http://foo.bar", head.Url)

	delete := c.Delete("http://foo.bar")
	as.eq("DELETE", delete.Method)
	as.eq("http://foo.bar", delete.Url)

	connect := c.Connect("http://foo.bar")
	as.eq("CONNECT", connect.Method)
	as.eq("http://foo.bar", connect.Url)

	options := c.Options("http://foo.bar")
	as.eq("OPTIONS", options.Method)
	as.eq("http://foo.bar", options.Url)

	trace := c.Trace("http://foo.bar")
	as.eq("TRACE", trace.Method)
	as.eq("http://foo.bar", trace.Url)
}
