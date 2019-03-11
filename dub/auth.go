package dub

import (
	"encoding/base64"
)

type AuthSpec interface {
	Token() string
}

type BasicAuth struct {
	User, Pass string
}

func (b *BasicAuth) Token() string {
	return "Basic " + base64.StdEncoding.EncodeToString(b.payload())
}

func (b *BasicAuth) payload() []byte {
	return []byte(b.User + ":" + b.Pass)
}

func NewBasicAuth(user, pass string) AuthSpec {
	return &BasicAuth{User: user, Pass: pass}
}

type TokenAuth struct {
	BearerToken string
}

func (t *TokenAuth) Token() string {
	return `Bearer ` + t.BearerToken
}

func NewTokenAuth(token string) AuthSpec {
	return &TokenAuth{BearerToken: token}
}
