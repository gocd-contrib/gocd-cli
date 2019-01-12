package dub

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type RawRequestHandler func(*http.Request) error
type ResponseBodyConsumer func(io.Reader) error
type ResponseHandler func(*Response) error
type ProgressHandler func(pr *Progress) error

type Opts struct {
	Cookies      []*http.Cookie
	Headers      map[string][]string
	Auth         AuthSpec
	ContentType  string
	OnProgress   []ProgressHandler
	OnBeforeSend []RawRequestHandler
}

var methodsCanHaveBody = map[string]struct{}{
	"put":    struct{}{},
	"patch":  struct{}{},
	"post":   struct{}{},
	"delete": struct{}{},
}

func allowBody(method string) bool {
	_, ok := methodsCanHaveBody[strings.ToLower(method)]
	return ok
}

type lengther interface {
	Len() int
}

type wrErr struct {
	cause  error
	reason string
}

func (w *wrErr) Error() string {
	return fmt.Sprintf("%s; cause:\n  %v", w.reason, w.cause)
}

func wrapErr(cause error, reason string) error {
	return &wrErr{cause: cause, reason: reason}
}
