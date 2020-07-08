package api

import (
	"errors"
	"fmt"
	"io"
	"path"
	"strconv"

	"github.com/gocd-contrib/gocd-cli/cfg"
	"github.com/gocd-contrib/gocd-cli/dub"
)

var (
	V1 = V(1)
	V2 = V(2)
	V3 = V(3)
	V4 = V(4)
	V5 = V(5)
	V6 = V(6)
	V7 = V(7)
)

type Builder struct {
	ApiVersion int

	conf *cfg.Config
	c    *dub.Client
}

func (b *Builder) Get(path string, onCreate ...CreateHook) *Req {
	return NewReq(b.c.Get(b.Url(path)), nil, b, onCreate)
}

func (b *Builder) Put(path string, data io.Reader, onCreate ...CreateHook) *Req {
	return NewReq(b.c.Put(b.Url(path)), data, b, b.AddConfirmHeaderIfBodyNil(onCreate))
}

func (b *Builder) Patch(path string, data io.Reader, onCreate ...CreateHook) *Req {
	return NewReq(b.c.Patch(b.Url(path)), data, b, b.AddConfirmHeaderIfBodyNil(onCreate))
}

func (b *Builder) Post(path string, data io.Reader, onCreate ...CreateHook) *Req {
	return NewReq(b.c.Post(b.Url(path)), data, b, b.AddConfirmHeaderIfBodyNil(onCreate))
}

func (b *Builder) Delete(path string, data io.Reader, onCreate ...CreateHook) *Req {
	return NewReq(b.c.Delete(b.Url(path)), data, b, onCreate)
}

func (b *Builder) AddConfirmHeaderIfBodyNil(onCreate []CreateHook) []CreateHook {
	return append(onCreate, func(req *dub.Request) error {
		if nil == req.Body {
			req.Header(`X-GoCD-Confirm`, `true`)
		}
		return nil
	})
}

func (b *Builder) Validate() error {
	return b.conf.WithBaseUrlValidation(b.conf.GetServerUrl(), nil)
}

func (b *Builder) Url(uri string) string {
	return b.conf.GetServerUrl() + path.Clean(uri)
}

func (b *Builder) AcceptHeader() string {
	return `application/vnd.go.cd.v` + strconv.Itoa(b.ApiVersion) + `+json`
}

func (b *Builder) Auth() (dub.AuthSpec, error) {
	auth := b.conf.GetAuth()

	if 0 == len(auth) {
		return nil, errors.New(`API auth is not configured`)
	}

	if err := checkAuth(auth, `type`); err != nil {
		return nil, err
	}

	switch auth[`type`] {
	case `basic`:
		if err := checkAuth(auth, `user`); err != nil {
			return nil, err
		}

		if err := checkAuth(auth, `password`); err != nil {
			return nil, err
		}

		return dub.NewBasicAuth(auth[`user`], auth[`password`]), nil
	case `token`:
		if err := checkAuth(auth, `token`); err != nil {
			return nil, err
		}

		return dub.NewTokenAuth(auth[`token`]), nil
	case `none`:
		return nil, nil
	default:
		return nil, fmt.Errorf(`Unknown authentication scheme: %q`, auth[`type`])
	}
}

func checkAuth(auth map[string]string, key string) error {
	if _, ok := auth[key]; !ok {
		return fmt.Errorf(`Failed to construct authentication spec; %q is missing`, key)
	}
	return nil
}

func V(version int) *Builder {
	return New(version, cfg.Conf(), dub.New())
}

func New(version int, config *cfg.Config, client *dub.Client) *Builder {
	return &Builder{ApiVersion: version, conf: config, c: client}
}
