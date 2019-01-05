package dub

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type mockRT func(*http.Request) (*http.Response, error)

func (f mockRT) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func okResp() *http.Response {
	return resp(200, "OK")
}

func nopRt() mockRT {
	return mockRT(func(rq *http.Request) (rs *http.Response, e error) {
		rs = okResp()
		return
	})
}

func nopCl() *Client {
	return Make(nopRt())
}

func testCl(f func(*http.Request) (*http.Response, error)) *Client {
	return Make(mockRT(f))
}

func resp(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       ioutil.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func ignoreResponse(res *Response) error {
	return res.Consume(func(r io.Reader) error {
		if r != nil {
			io.Copy(devnull(), r)
		}
		return nil
	})
}

func reqfmt(s ...string) string {
	return strings.Join(s, "\r\n")
}

type drain struct{}

func (d *drain) Write(b []byte) (int, error) { return len(b), nil }
func (d *drain) Close() error                { return nil }

func devnull() io.WriteCloser {
	return &drain{}
}

type dummyAuth struct {
	s string
}

func (d *dummyAuth) Token() string {
	return "Dummy " + d.s
}

type asserter struct {
	t *testing.T
}

func fakeAuth(s string) AuthSpec {
	return &dummyAuth{s: s}
}

func (a *asserter) eq(expected, actual interface{}) {
	a.t.Helper()
	if expected != actual {
		a.t.Errorf("Expected %v to equal %v", actual, expected)
	}
}

func (a *asserter) neq(expected, actual interface{}) {
	a.t.Helper()
	if expected == actual {
		a.t.Errorf("Expected %v to not equal %v", actual, expected)
	}
}

func (a *asserter) deepEqI64(expected, actual []int64) {
	a.t.Helper()

	if !func() bool {
		if len(expected) != len(actual) {
			return false
		}

		for i, v := range actual {
			if expected[i] != v {
				return false
			}
		}

		return true
	}() {
		a.t.Errorf("Expected %v to equal %v", actual, expected)
	}
}

func (a *asserter) isNil(actual interface{}) {
	a.t.Helper()
	if nil != actual {
		a.t.Errorf("Expected %v to be nil", actual)
	}
}

func (a *asserter) err(expected string, e error) {
	a.t.Helper()
	if nil == e {
		a.t.Errorf("Expected error %q, but got nil", expected)
		return
	}

	if e.Error() != expected {
		a.t.Errorf("Expected error %q, but got %q", expected, e)
	}
}

func (a *asserter) ok(err error) {
	a.t.Helper()
	if nil != err {
		a.t.Errorf("Expected no error, but got %v", err)
	}
}

func (a *asserter) is(b bool) {
	a.t.Helper()
	if !b {
		a.t.Errorf("Expected to be true")
	}
}

func (a *asserter) not(b bool) {
	a.t.Helper()
	if b {
		a.t.Errorf("Expected to be false")
	}
}

func asserts(t *testing.T) *asserter {
	return &asserter{t: t}
}
