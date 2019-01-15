package api_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gocd-contrib/gocd-cli/api"
	"github.com/gocd-contrib/gocd-cli/cfg"
	"github.com/gocd-contrib/gocd-cli/dub"
	"github.com/spf13/afero"
)

const (
	TEST_CONF_FILE = `conf.yaml`
)

func testApi(version int, conf *cfg.Config) *api.Builder {
	return api.New(version, conf, testDub())
}

func testDub() *dub.Client {
	return dub.Make(nopRt())
}

func testConf() *cfg.Config {
	c := cfg.NewConfig(afero.NewMemMapFs())
	c.Consume("")
	return c
}

func writeContent(fs afero.Fs, path, content string) error {
	if f, err := fs.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.FileMode(0644)); err != nil {
		return err
	} else {
		defer f.Close()

		_, err = f.WriteString(content)
		return err
	}
}

func makeConf(content string) (*cfg.Config, error) {
	fs := afero.NewMemMapFs()
	c := cfg.NewConfig(fs)

	if err := writeContent(fs, TEST_CONF_FILE, content); err != nil {
		return nil, err
	} else {
		if err = c.Consume(TEST_CONF_FILE); err != nil {
			return nil, err
		}
	}

	return c, nil
}

type mockRT func(*http.Request) (*http.Response, error)

func (f mockRT) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func okResp() *http.Response {
	return resp(200, `OK`)
}

func nopRt() mockRT {
	return mockRT(func(rq *http.Request) (rs *http.Response, e error) {
		rs = okResp()
		return
	})
}

func resp(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       ioutil.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

type dummyAuth struct {
	s string
}

func (d *dummyAuth) Token() string {
	return "Dummy " + d.s
}

type testConfigurer struct{}

func (tc *testConfigurer) AcceptHeader() string {
	return `accept/me`
}

func (tc *testConfigurer) Auth() (dub.AuthSpec, error) {
	return &dummyAuth{s: `token`}, nil
}

type asserter struct {
	t *testing.T
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
