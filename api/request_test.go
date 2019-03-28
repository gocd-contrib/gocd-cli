package api_test

import (
	"errors"
	"testing"

	"github.com/gocd-contrib/gocd-cli/api"
	"github.com/gocd-contrib/gocd-cli/dub"
)

const (
	TEST_API_URL = `http://test/go/api/doit`
)

func TestRequestValidateUrl(t *testing.T) {
	as := asserts(t)

	req := &api.Req{
		Raw:        testDub().Get(TEST_API_URL),
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
		Raw:        testDub().Get(TEST_API_URL),
		Configurer: &testConfigurer{},
	}

	as.ok(req.Config())

	as.eq(`accept/me`, req.Raw.Headers.Get(`Accept`))
	as.eq(`Dummy token`, req.Raw.Headers.Get(`Authorization`))
}

func TestConfigAppliesOnCreateHooks(t *testing.T) {
	as := asserts(t)

	req := &api.Req{
		Raw:        testDub().Get(TEST_API_URL),
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

func TestSendWithSuccessResponse(t *testing.T) {
	as := asserts(t)

	req := &api.Req{
		Raw:        testDub().Get(TEST_API_URL),
		Configurer: &testConfigurer{},
	}

	didRun := false
	onSuccess := func(res *dub.Response) error {
		as.is(res.IsSuccess())
		didRun = true
		return nil
	}

	onFail := func(res *dub.Response) error {
		t.Error(`Response should not be failure!`)
		return nil
	}

	as.ok(req.Send(onSuccess, onFail))
	as.is(didRun)
}

func TestSendWithFailureResponse(t *testing.T) {
	as := asserts(t)

	req := &api.Req{
		Raw:        dub.Make(useRt(404, `not found`)).Get(TEST_API_URL),
		Configurer: &testConfigurer{},
	}

	didRun := false
	onFail := func(res *dub.Response) error {
		as.is(res.IsError())
		didRun = true
		return nil
	}

	onSuccess := func(res *dub.Response) error {
		t.Error(`Response should not be success!`)
		return nil
	}

	as.ok(req.Send(onSuccess, onFail))
	as.is(didRun)
}

func TestAbortsSendByReturningErrorOnCreate(t *testing.T) {
	as := asserts(t)

	req := &api.Req{
		Raw:        testDub().Get(TEST_API_URL),
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

	as.err(`Stop this request!`, req.Send(onResponse, onResponse))
	as.not(didRun)
}
