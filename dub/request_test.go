package dub

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestRequestConfigByOpts(t *testing.T) {
	as := asserts(t)

	didRun := false

	c := testCl(func(r *http.Request) (*http.Response, error) {
		as.eq("http://test", r.URL.String())
		as.eq(2, len(r.Cookies()))
		as.eq(4, len(r.Header)) // Auth, ContType, X-Foo, Cookie
		as.eq("Dummy secret", r.Header.Get("Authorization"))
		as.eq("text/html", r.Header.Get("Content-Type"))
		as.eq("bar", r.Header.Get("X-Foo"))

		didRun = true
		return okResp(), nil
	})

	req := c.Get("http://test").Opts(&Opts{
		Cookies: []*http.Cookie{
			&http.Cookie{Name: `foo`, Value: `bar`},
			&http.Cookie{Name: `baz`, Value: `quu`},
		},
		Headers: map[string][]string{
			"X-Foo": {"bar"},
		},
		Auth:         fakeAuth("secret"),
		ContentType:  "text/html",
		OnProgress:   []ProgressHandler{func(p *Progress) error { return nil }},
		OnBeforeSend: []RawRequestHandler{func(p *http.Request) error { return nil }},
	})

	as.neq(nil, req)
	as.eq("http://test", req.Url)
	as.eq(3, len(req.Headers))
	as.eq("Dummy secret", req.Headers.Get("Authorization"))
	as.eq("text/html", req.Headers.Get("Content-Type"))
	as.eq("bar", req.Headers.Get("X-Foo"))
	as.eq(2, len(req.Cookies))
	as.eq(1, len(req.onProgress))
	as.eq(1, len(req.onBeforeSend))

	as.ok(req.Do(ignoreResponse))
	as.is(didRun)
}

func TestRequestUrl(t *testing.T) {
	as := asserts(t)

	didRun := false

	c := testCl(func(r *http.Request) (*http.Response, error) {
		as.eq("http://test", r.URL.String())

		didRun = true
		return okResp(), nil
	})

	req := c.Get("http://test")

	as.neq(nil, req)
	as.eq("http://test", req.Url)

	as.ok(req.Do(ignoreResponse))
	as.is(didRun)
}

func TestWhichRequestsAcceptBody(t *testing.T) {
	as := asserts(t)

	c := nopCl()
	url := "http://test"

	bodyDeny := []*Request{c.Get(url), c.Head(url), c.Connect(url), c.Options(url), c.Trace(url)}
	for _, req := range bodyDeny {
		as.err(fmt.Sprintf("Method `%s` does not accept a request body", req.Method), req.DataString("test").Do(ignoreResponse))
	}

	bodyAccept := []*Request{c.Delete(url), c.Post(url), c.Put(url), c.Patch(url)}
	for _, req := range bodyAccept {
		as.ok(req.DataString("test").Do(ignoreResponse))
	}
}

func TestRequestBeforeSendAllowsNativeRequestAccess(t *testing.T) {
	as := asserts(t)

	c := nopCl()

	didRun := false
	var native *http.Request

	as.ok(c.Get("http://test").BeforeSend(func(r *http.Request) error {
		native = r
		didRun = true
		return nil
	}).Do(ignoreResponse))

	as.is(didRun)
	as.neq(nil, native)
}

func TestRequestData(t *testing.T) {
	as := asserts(t)

	var actual string
	didRun := false

	c := testCl(func(r *http.Request) (*http.Response, error) {
		b := &strings.Builder{}

		as.neq(nil, r.Body)

		_, err := io.Copy(b, r.Body)
		as.ok(err)

		actual = b.String()

		as.eq(int64(5), r.ContentLength)
		didRun = true
		return okResp(), nil
	})

	as.ok(c.Put("http://test").Data(strings.NewReader("hello")).Do(ignoreResponse))
	as.is(didRun)
	as.eq("hello", actual)
}

func TestRequestDataString(t *testing.T) {
	as := asserts(t)

	var actual string
	didRun := false

	c := testCl(func(r *http.Request) (*http.Response, error) {
		b := &strings.Builder{}

		as.neq(nil, r.Body)
		_, err := io.Copy(b, r.Body)
		as.ok(err)
		actual = b.String()

		as.eq(int64(5), r.ContentLength)
		didRun = true
		return okResp(), nil
	})

	as.ok(c.Put("http://test").DataString("hello").Do(ignoreResponse))
	as.is(didRun)
	as.eq("hello", actual)
}

func TestRequestAuth(t *testing.T) {
	as := asserts(t)

	didRun := false

	c := testCl(func(r *http.Request) (*http.Response, error) {
		as.eq("Dummy secret", r.Header.Get("Authorization"))

		didRun = true
		return okResp(), nil
	})

	as.ok(c.Get("http://test").Auth(fakeAuth("secret")).Do(ignoreResponse))
	as.is(didRun)
}

func TestRequestContentType(t *testing.T) {
	as := asserts(t)

	didRun := false

	c := testCl(func(r *http.Request) (*http.Response, error) {
		as.eq("application/json", r.Header.Get("Content-Type"))

		didRun = true
		return okResp(), nil
	})

	as.ok(c.Get("http://test").ContentType("application/json").Do(ignoreResponse))
	as.is(didRun)
}

func TestAddQuery(t *testing.T) {
	as := asserts(t)

	req := &Request{
		Url: "http://test",
	}

	req.AddQuery(map[string][]string{
		`hello`: {`world`, `monde`},
		`foo`:   {`bar`},
	})

	as.eq(`http://test?foo=bar&hello=world&hello=monde`, req.Url)
}

func TestAddQueryToExistingQueryString(t *testing.T) {
	as := asserts(t)

	req := &Request{
		Url: "http://test?foo=bar",
	}

	req.AddQuery(map[string][]string{
		`hello`: {`world`},
		`foo`:   {`baz`},
	})

	as.eq(`http://test?foo=bar&foo=baz&hello=world`, req.Url)
}

func TestRequestHeader(t *testing.T) {
	as := asserts(t)

	didRun := false

	c := testCl(func(r *http.Request) (*http.Response, error) {
		as.eq("Bar", r.Header.Get("Foo"))

		didRun = true
		return okResp(), nil
	})

	as.ok(c.Get("http://test").Header("Foo", "Bar").Do(ignoreResponse))
	as.is(didRun)
}

func TestRequestCookie(t *testing.T) {
	as := asserts(t)

	didRun := false

	c := testCl(func(r *http.Request) (*http.Response, error) {
		as.eq(1, len(r.Cookies()))
		ck, err := r.Cookie(`cookie-monster`)

		as.ok(err)

		as.neq(nil, ck)
		as.eq(`cookie-monster=me-want-cookie`, ck.String())

		didRun = true
		return okResp(), nil
	})

	as.ok(c.Get("http://test").Cookie(&http.Cookie{
		Name:  "cookie-monster",
		Value: "me-want-cookie",
	}).Do(ignoreResponse))
	as.is(didRun)
}

func TestRequestSetHeaders(t *testing.T) {
	as := asserts(t)

	didRun := false

	c := testCl(func(r *http.Request) (*http.Response, error) {
		as.eq("Bar", r.Header.Get("Foo"))
		as.eq("Bye", r.Header.Get("Hi"))

		w := &strings.Builder{}
		r.Header.Write(w) // headers can repeat (multi-value)
		as.eq(reqfmt(
			`Foo: Bar`,
			`Foo: Baz`,
			`Hi: Bye`,
			``,
		), w.String())

		didRun = true
		return okResp(), nil
	})

	as.ok(c.Get("http://test").SetHeaders(map[string][]string{
		"Foo": {"Bar", "Baz"},
		"Hi":  {"Bye"},
	}).Do(ignoreResponse))
	as.is(didRun)
}
