package api

import (
	"io"

	"github.com/gocd-contrib/gocd-cli/dub"
	"github.com/gocd-contrib/gocd-cli/utils"
)

type debugWriter func([]byte) error

func (dw debugWriter) Write(b []byte) (int, error) {
	return len(b), dw(b)
}

func traceRequest(r *dub.Request) {
	utils.Debug(`Sending API request %s %s`, r.Method, r.Url)

	if utils.DebugMode {
		utils.Debug(`Headers >>>`)
		for h, vals := range r.Headers {
			if "Authorization" == h {
				vals = []string{`:: REDACTED ::`}
			}

			for _, v := range vals {
				utils.Debug(`%s: %s`, h, v)
			}
		}

		if r.Body != nil {
			bodyTap := debugWriter(func(b []byte) error {
				utils.Debug(string(b))
				return nil
			})

			utils.Debug(`Body >>>`)

			if mult, ok := r.Body.(*dub.Multipart); ok {
				mult.MultipartPayload = dub.NewWireTapPayload(mult.MultipartPayload, bodyTap)
			} else {
				r.Body = io.TeeReader(r.Body, bodyTap)
			}
		}
	}
}
