package dub

import (
	"io"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Status  int
	Headers http.Header
	Raw     *http.Response

	onProgress []ProgressHandler
}

func (r *Response) IsAuthError() bool {
	return r.Status == 401
}

func (r *Response) IsError() bool {
	return r.Status >= 400
}

func (r *Response) IsSuccessOrRedirect() bool {
	return r.Status < 400 && r.Status > 199
}

func (r *Response) IsSuccess() bool {
	return r.Status < 300 && r.Status > 199
}

func (r *Response) IsRedirect() bool {
	return r.Status < 400 && r.Status > 299
}

func (r *Response) OnProgress(handler ProgressHandler) *Response {
	r.onProgress = append(r.onProgress, handler)
	return r
}

func (r *Response) Consume(doRead ResponseBodyConsumer) error {
	defer r.Raw.Body.Close()

	if len(r.onProgress) > 0 {
		progress := newProgress(r.Raw.ContentLength, r.onProgress)
		progress.RawResponse = r.Raw
		return doRead(io.TeeReader(r.Raw.Body, newProgressWriter(progress)))
	} else {
		return doRead(r.Raw.Body)
	}
}

func (r *Response) ReadAll() ([]byte, error) {
	var payload []byte
	if err := r.Consume(func(body io.Reader) (err error) {
		payload, err = ioutil.ReadAll(body)
		return err
	}); err != nil {
		return nil, err
	}
	return payload, nil
}

func newResp(r *http.Response) *Response {
	return &Response{Raw: r, Headers: r.Header, Status: r.StatusCode}
}
