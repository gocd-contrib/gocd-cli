package api

import (
	"io"
	"io/ioutil"

	"github.com/gocd-contrib/gocd-cli/dub"
	"github.com/gocd-contrib/gocd-cli/utils"
)

func ReadBodyAndDo(res *dub.Response, action func([]byte) error) error {
	return res.Consume(func(reader io.Reader) error {
		if b, err := ioutil.ReadAll(reader); err != nil {
			return utils.InspectError(err, `reading preflight response from %q`, res.Raw.Request.URL)
		} else {
			return action(b)
		}
	})
}
