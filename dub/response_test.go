package dub

import (
	"io"
	"strings"
	"testing"
)

func TestResponseStatus(t *testing.T) {
	as := asserts(t)

	as.is((&Response{Status: 200}).IsSuccess())
	as.is((&Response{Status: 204}).IsSuccess())
	as.is((&Response{Status: 204}).IsSuccessOrRedirect())
	as.not((&Response{Status: 201}).IsRedirect())
	as.not((&Response{Status: 200}).IsError())

	as.is((&Response{Status: 302}).IsRedirect())
	as.is((&Response{Status: 307}).IsSuccessOrRedirect())
	as.not((&Response{Status: 300}).IsSuccess())
	as.not((&Response{Status: 300}).IsError())

	as.is((&Response{Status: 404}).IsError())
	as.is((&Response{Status: 500}).IsError())
	as.not((&Response{Status: 401}).IsSuccessOrRedirect())
	as.not((&Response{Status: 403}).IsRedirect())
	as.not((&Response{Status: 422}).IsSuccess())
}

func TestResponseConsume(t *testing.T) {
	as := asserts(t)

	r := newResp(resp(200, "hello"))

	as.neq(nil, r.Raw)
	as.neq(nil, r.Raw.Body)

	var body string
	as.ok(r.Consume(func(re io.Reader) (err error) {
		b := &strings.Builder{}
		if _, err = io.Copy(b, re); err == nil {
			body = b.String()
		}
		return
	}))

	as.eq("hello", body)
}

func TestResponseReadAll(t *testing.T) {
	as := asserts(t)

	r := newResp(resp(200, "hello"))

	as.neq(nil, r.Raw)
	as.neq(nil, r.Raw.Body)

	body, err := r.ReadAll()
	as.ok(err)
	as.eq("hello", string(body))
}
