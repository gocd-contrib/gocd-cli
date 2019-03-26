package api

import (
	"errors"
	"io"
	"net/url"

	"github.com/gocd-contrib/gocd-cli/dub"
	"github.com/gocd-contrib/gocd-cli/utils"
)

type RequestConfigurer interface {
	AcceptHeader() string
	Auth() (dub.AuthSpec, error)
	Validate() error
}

type CreateHook func(*dub.Request) error

type Req struct {
	Raw        *dub.Request
	OnCreate   []CreateHook
	Configurer RequestConfigurer
}

func (r *Req) Send(onResponse, onErrorResponse func(*dub.Response) error) error {
	if err := r.Config(); err != nil {
		return utils.InspectError(err, `configuring api.Req during Send()`)
	} else {
		handleResponse := func(res *dub.Response) error {
			if res.IsError() {
				utils.Debug(`handling error response %d`, res.Status)
				if onErrorResponse != nil {
					return onErrorResponse(res)
				}
				utils.Debug(`error response handler is nil; ignoring response`)
				return nil
			}
			utils.Debug(`handling success response %d`, res.Status)
			if onResponse != nil {
				return onResponse(res)
			}
			utils.Debug(`success response handler is nil; ignoring response`)
			return nil
		}

		traceRequest(r.Raw)
		return utils.InspectError(r.Raw.Do(handleResponse), `making api %s request to %q`, r.Raw.Method, r.Raw.Url)
	}
}

func (r *Req) ValidateUrl() error {
	if err := r.Configurer.Validate(); err != nil {
		return utils.InspectError(err, `validating API request builder configuration`)
	}

	if u, err := url.Parse(r.Raw.Url); err == nil {
		if "" == u.Scheme || "" == u.Host {
			return errors.New("API URL is not absolute; make sure you have configured `server-url`")
		}
	} else {
		return utils.InspectError(err, `parsing URL %q during ValidateUrl()`, r.Raw.Url)
	}
	return nil
}

func (r *Req) Config() error {
	rc := r.Configurer

	if err := r.ValidateUrl(); err != nil {
		return utils.InspectError(err, `validating url %q`, r.Raw.Url)
	}

	r.Raw.Header(`Accept`, rc.AcceptHeader())

	if auth, err := rc.Auth(); err == nil {
		r.Raw.Auth(auth)
	} else {
		return utils.InspectError(err, `setting auth on api.Req`)
	}

	if len(r.OnCreate) > 0 {
		for _, hook := range r.OnCreate {
			if err := hook(r.Raw); err != nil {
				return utils.InspectError(err, `executing api.Req onCreate hook %v`, hook)
			}
		}
	}

	return nil
}

func NewReq(raw *dub.Request, data io.Reader, rc RequestConfigurer, createHooks []CreateHook) *Req {
	raw.Data(data)
	return &Req{
		Raw:        raw,
		OnCreate:   createHooks,
		Configurer: rc,
	}
}
