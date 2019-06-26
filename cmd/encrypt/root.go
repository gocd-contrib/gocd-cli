package encrypt

import (
	"encoding/json"

	"github.com/gocd-contrib/gocd-cli/api"
	"github.com/gocd-contrib/gocd-cli/dub"
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a plaintext value with your GoCD server's secrets",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		root.Run(args)
	},
}

var root = &rootRunner{}

type rootRunner struct{}

func (r *rootRunner) Run(args []string) {
	api.V1.Post(`/api/admin/encrypt`, api.JsonReader(map[string]string{
		`value`: args[0],
	}), addContentType).Send(r.onSuccess, r.onFail)
}

func (r *rootRunner) onSuccess(res *dub.Response) error {
	return api.ReadBodyAndDo(res, func(b []byte) error {
		v := &encryptedValue{}
		if err := json.Unmarshal(b, v); err == nil {
			utils.Echof(v.Value)
			return nil
		} else {
			return utils.InspectError(err, `parsing encrypted-value response`)
		}
	})
}

func (r *rootRunner) onFail(res *dub.Response) error {
	return api.ReadBodyAndDo(res, func(b []byte) error {
		api.DieOnAuthError(res)
		api.DieOn4XX(res, b, api.ParseMessage)
		api.DieOnUnexpected(res, b)

		return nil
	})
}

type encryptedValue struct {
	Value string `json:"encrypted_value"`
}

func addContentType(req *dub.Request) error {
	req.Header(`Content-Type`, `application/json`)
	return nil
}
