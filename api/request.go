package api

import (
	"errors"
	"io"
	"net/url"

	"github.com/gocd-contrib/gocd-cli/dub"
)

type RequestConfigurer interface {
	AcceptHeader() string
	Auth() (dub.AuthSpec, error)
}

type CreateHook func(*dub.Request) error

type Req struct {
	Raw        *dub.Request
	OnCreate   []CreateHook
	Configurer RequestConfigurer
}

func (r *Req) Send(onResponse func(*dub.Response) error) error {
	if err := r.Config(); err != nil {
		return err
	} else {
		return r.Raw.Do(onResponse)
	}
}

func (r *Req) ValidateUrl() error {
	if u, err := url.Parse(r.Raw.Url); err == nil {
		if "" == u.Scheme || "" == u.Host {
			return errors.New("API URL is not absolute; make sure you have configured `server-url`")
		}
	} else {
		return err
	}
	return nil
}

func (r *Req) Config() error {
	rc := r.Configurer

	if err := r.ValidateUrl(); err != nil {
		return err
	}

	r.Raw.Header(`Accept`, rc.AcceptHeader())

	if auth, err := rc.Auth(); err == nil {
		r.Raw.Auth(auth)
	} else {
		return err
	}

	if len(r.OnCreate) > 0 {
		for _, hook := range r.OnCreate {
			if err := hook(r.Raw); err != nil {
				return err
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
