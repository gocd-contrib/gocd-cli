package api_test

import (
	"io"
	"strings"
	"testing"

	"github.com/gocd-contrib/gocd-cli/api"
)

func TestAcceptHeader(t *testing.T) {
	as := asserts(t)
	as.eq(`application/vnd.go.cd.v999+json`, testApi(999, testConf()).AcceptHeader())
}

func TestAuthFailsWithoutConfiguration(t *testing.T) {
	as := asserts(t)
	v := testApi(999, testConf())

	_, err := v.Auth()
	as.err(`API auth is not configured`, err)
}

func TestAuthFailsWhenUnknownAuthType(t *testing.T) {
	as := asserts(t)

	c, err := makeConf(`
auth:
  type: foo
`)
	as.ok(err)
	v := testApi(999, c)
	_, err = v.Auth()
	as.err(`Unknown authentication scheme: "foo"`, err)
}

func TestAuthFailsWhenMissingType(t *testing.T) {
	as := asserts(t)

	c, err := makeConf(`
auth:
  hash: wha?
`)
	as.ok(err)
	v := testApi(999, c)
	_, err = v.Auth()
	as.err(`Failed to construct authentication spec; "type" is missing`, err)
}
func TestAuthFailsWhenMissingUser(t *testing.T) {
	as := asserts(t)

	c, err := makeConf(`
auth:
  type: basic
  password: wha?
`)
	as.ok(err)
	v := testApi(999, c)
	_, err = v.Auth()
	as.err(`Failed to construct authentication spec; "user" is missing`, err)
}

func TestAuthFailsWhenMissingPassword(t *testing.T) {
	as := asserts(t)

	c, err := makeConf(`
auth:
  type: basic
  user: foo
`)
	as.ok(err)
	v := testApi(999, c)
	_, err = v.Auth()
	as.err(`Failed to construct authentication spec; "password" is missing`, err)
}

func TestMethods(t *testing.T) {
	as := asserts(t)
	v := testApi(999, testConf())

	// GET has a different method signature
	cases := map[string]func(string, io.Reader, ...api.CreateHook) *api.Req{
		`POST`:   v.Post,
		`PUT`:    v.Put,
		`PATCH`:  v.Patch,
		`DELETE`: v.Delete,
	}

	for method, factory := range cases {
		req := factory(`/api/path`, nil)
		as.eq(method, req.Raw.Method)
	}

	// also test GET
	as.eq(`GET`, v.Get(`/api/path`).Raw.Method)
}

func TestUpdateMethodsAddConfirmHeader(t *testing.T) {
	as := asserts(t)

	v := testApi(999, testConf())

	cases := []func(string, io.Reader, ...api.CreateHook) *api.Req{
		v.Post,
		v.Put,
		v.Patch,
	}

	for _, factory := range cases {
		req := factory(`/api/path`, nil)
		as.eq(1, len(req.OnCreate))
		req.OnCreate[0](req.Raw)
		as.eq(`true`, req.Raw.Headers.Get(`X-GoCD-Confirm`))
	}

	for _, factory := range cases {
		req := factory(`/api/path`, strings.NewReader(`should not set confirm header when body is present`))
		as.eq(1, len(req.OnCreate))
		req.OnCreate[0](req.Raw)
		_, isSet := req.Raw.Headers[`X-GoCD-Confirm`]
		as.not(isSet)
	}
}

func TestIntegrationUpdateMethodsHaveConfirmHeaderUponConfig(t *testing.T) {
	as := asserts(t)

	c := testConf()
	as.ok(c.SetServerUrl(`http://test/go`))
	as.ok(c.SetBasicAuth(`foo`, `bar`))

	v := testApi(999, c)

	cases := []func(string, io.Reader, ...api.CreateHook) *api.Req{
		v.Post,
		v.Put,
		v.Patch,
	}

	for _, factory := range cases {
		req := factory(`/api/path`, nil)
		as.ok(req.Config())
		as.eq(`true`, req.Raw.Headers.Get(`X-GoCD-Confirm`))
	}

	for _, factory := range cases {
		req := factory(`/api/path`, strings.NewReader(`should not set confirm header when body is present`))
		as.ok(req.Config())
		_, isSet := req.Raw.Headers[`X-GoCD-Confirm`]
		as.not(isSet)
	}
}
