package api_test

import (
	"errors"
	"testing"

	"github.com/gocd-contrib/gocd-cli/api"
	"github.com/gocd-contrib/gocd-cli/dub"
)

func TestRequestValidateUrl(t *testing.T) {
	as := asserts(t)

	req := &api.Req{
		Raw:        testDub().Get(`http://test`),
		Configurer: &testConfigurer{},
	}

	as.ok(req.ValidateUrl())

	req = &api.Req{
		Raw:        testDub().Get(`/foo/bar`),
		Configurer: &testConfigurer{},
	}

	as.err("API URL is not absolute; make sure you have configured `server-url`", req.ValidateUrl())
}

func TestConfigSetsAcceptHeaderAndAuth(t *testing.T) {
	as := asserts(t)

	req := &api.Req{
		Raw:        testDub().Get(`http://test`),
		Configurer: &testConfigurer{},
	}

	as.ok(req.Config())

	as.eq(`accept/me`, req.Raw.Headers.Get(`Accept`))
	as.eq(`Dummy token`, req.Raw.Headers.Get(`Authorization`))
}

func TestConfigAppliesOnCreateHooks(t *testing.T) {
	as := asserts(t)

	req := &api.Req{
		Raw:        testDub().Get(`http://test`),
		Configurer: &testConfigurer{},
		OnCreate: []api.CreateHook{
			func(r *dub.Request) error {
				r.Header(`Foo`, `bar`)
				return nil
			},
		},
	}

	as.ok(req.Config())

	as.eq(`bar`, req.Raw.Headers.Get(`Foo`))
}

func TestSend(t *testing.T) {
	as := asserts(t)

	req := &api.Req{
		Raw:        testDub().Get(`http://test`),
		Configurer: &testConfigurer{},
	}

	didRun := false
	onResponse := func(res *dub.Response) error {
		as.is(res.IsSuccess())
		didRun = true
		return nil
	}

	as.ok(req.Send(onResponse))
	as.is(didRun)
}

func TestAbortsSendByReturningErrorOnCreate(t *testing.T) {
	as := asserts(t)

	req := &api.Req{
		Raw:        testDub().Get(`http://test`),
		Configurer: &testConfigurer{},
		OnCreate: []api.CreateHook{
			func(r *dub.Request) error {
				return errors.New(`Stop this request!`)
			},
		},
	}

	didRun := false
	onResponse := func(res *dub.Response) error {
		t.Error(`Request should not be sent!`)
		return nil
	}

	as.err(`Stop this request!`, req.Send(onResponse))
	as.not(didRun)
}
