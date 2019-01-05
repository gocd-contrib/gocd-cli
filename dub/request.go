package dub

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Request struct {
	Url, Method string

	Headers http.Header
	Body    io.Reader
	Raw     *http.Request

	onBeforeSend []RequestHandler
	onProgress   []ProgressHandler
	c            *Client
}

func (r *Request) Opts(opts *Opts) *Request {
	if len(opts.Headers) > 0 {
		r.SetHeaders(opts.Headers)
	}

	if opts.Auth != nil {
		r.Auth(opts.Auth)
	}

	if opts.ContentType != "" {
		r.ContentType(opts.ContentType)
	}

	if opts.OnProgress != nil {
		r.onProgress = opts.OnProgress
	}

	if opts.OnBeforeSend != nil {
		r.onBeforeSend = opts.OnBeforeSend
	}

	return r
}

func (r *Request) OnProgress(handler ProgressHandler) *Request {
	r.onProgress = append(r.onProgress, handler)
	return r
}

func (r *Request) SetHeaders(headers map[string][]string) *Request {
	r.Headers = make(http.Header, len(headers))
	if len(headers) > 0 {
		for key, val := range headers {
			r.Headers[http.CanonicalHeaderKey(key)] = val
		}
	}
	return r
}

func (r *Request) Header(key, value string) *Request {
	r.ensureHeaders().Add(key, value)
	return r
}

func (r *Request) ContentType(contentType string) *Request {
	r.ensureHeaders().Set("Content-Type", contentType)
	return r
}

func (r *Request) bodySize() int64 {
	switch v := r.Body.(type) {
	case lengther:
		if size := v.Len(); size > -1 {
			return int64(size)
		}
	case *os.File:
		if i, err := v.Stat(); err == nil {
			return i.Size()
		}
	}
	return int64(-1)
}

func (r *Request) setContentLength(req *http.Request) {
	if r.Body == nil {
		return
	}

	req.ContentLength = r.bodySize()
}

func (r *Request) ensureHeaders() http.Header {
	if nil == r.Headers {
		r.Headers = make(http.Header)
	}
	return r.Headers
}

func (r *Request) Auth(auth AuthSpec) *Request {
	if nil == r.Headers {
		r.Headers = make(http.Header)
	}

	if nil == auth {
		r.ensureHeaders().Del("Authorization")
	} else {
		r.ensureHeaders().Set("Authorization", auth.Token())
	}

	return r
}

func (r *Request) Data(data io.Reader) *Request {
	r.Body = data

	if m, ok := data.(*Multipart); ok {
		r.ContentType(m.ContentType())
	}
	return r
}

func (r *Request) DataString(data string) *Request {
	r.Data(strings.NewReader(data))
	return r
}

func (r *Request) BeforeSend(handler RequestHandler) *Request {
	r.onBeforeSend = append(r.onBeforeSend, handler)
	return r
}

func (r *Request) Do(onResponse ResponseHandler) error {
	if !allowBody(r.Method) && r.Body != nil {
		return fmt.Errorf("Method `%s` does not accept a request body", r.Method)
	}

	var body io.Reader
	var progress *Progress

	if r.Body != nil {
		if m, ok := r.Body.(*Multipart); ok {
			if err := m.Assemble(); err != nil {
				return wrapErr(err, "Failed to assemble multipart request body stream")
			}
		}

		if len(r.onProgress) > 0 {
			progress = newProgress(r.bodySize(), r.onProgress)
			body = io.TeeReader(r.Body, newProgressWriter(progress))
		} else {
			body = r.Body
		}
	}

	if req, errRq := http.NewRequest(r.Method, r.Url, body); errRq == nil {
		r.Raw = req

		if progress != nil {
			progress.RawRequest = req
		}

		req.Header = r.Headers
		r.setContentLength(req)

		if len(r.onBeforeSend) > 0 {
			for _, h := range r.onBeforeSend {
				if err := h(r); err != nil {
					return wrapErr(err, "Request.BeforeSend() hook failed")
				}
			}
		}

		res, errRs := r.c.native.Do(req)

		if errRs != nil {
			return errRs
		}

		return onResponse(newResp(res))
	} else {
		return wrapErr(errRq, "Failed to build native *http.Request")
	}
}
