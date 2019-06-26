package api

import (
	"io"
	"io/ioutil"

	"github.com/gocd-contrib/gocd-cli/dub"
	"github.com/gocd-contrib/gocd-cli/utils"
)

func ReadBodyAndDo(res *dub.Response, action func([]byte) error) error {
	if utils.DebugMode {
		utils.Debug(`Response status code: %d`, res.Status)
		utils.Debug(`Response Headers >>>`)

		for h, vals := range res.Headers {
			for _, v := range vals {
				utils.Debug(`%s: %s`, h, v)
			}
		}
	}

	return res.Consume(func(reader io.Reader) error {
		if b, err := ioutil.ReadAll(reader); err != nil {
			return utils.InspectError(err, `reading response from %q`, res.Raw.Request.URL)
		} else {
			utils.Debug("Response Body >>>\n%s", string(b))
			return action(b)
		}
	})
}

func DieOnUnexpected(res *dub.Response, body []byte) {
	if 500 <= res.Status {
		utils.DieLoudly(1, `Server responded with %d: %s`, res.Status, string(body))
	}
}

func DieOn4XX(res *dub.Response, body []byte, bodyParser func([]byte) (MessageResponse, error)) {
	if 399 < res.Status && 500 > res.Status {
		if bodyParser == nil {
			bodyParser = ParseMessage
		}
		if m, err := bodyParser(body); err == nil {
			utils.DieLoudly(1, m.String())
		} else {
			utils.InspectError(err, `parsing %d response message: %q`, res.Status, string(body))
			utils.DieLoudly(1, `Server responded with %d: %s`, res.Status, string(body))
		}
	}
}

func DieOnEtagStale(res *dub.Response, errorMsg string, t ...interface{}) {
	if 412 == res.Status {
		utils.DieLoudly(1, errorMsg, t...)
	}
}

func DieOnNotFound(res *dub.Response, errorMsg string, t ...interface{}) {
	if res.IsNotFound() {
		utils.DieLoudly(1, errorMsg, t...)
	}
}

func DieOnAuthError(res *dub.Response) {
	if res.IsAuthError() {
		utils.DieLoudly(1, `Invalid credentials. Either the configured username, password, or auth token is incorrect`)
	}
}
