package configrepo

import (
	"net/url"
	"path"

	"github.com/gocd-contrib/gocd-cli/api"
	"github.com/gocd-contrib/gocd-cli/dub"
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var RmCmd = &cobra.Command{
	Use:   "rm id",
	Short: "Deletes a config-repo by id",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rm.Run(args)
	},
}

var rm = &RmRunner{}

type RmRunner struct{}

func (r *RmRunner) Run(args []string) {
	if err := api.V1.Delete(r.url(args[0]), nil).Send(r.onSuccess, r.onFail); err != nil {
		utils.AbortLoudly(err)
	}
}

func (r *RmRunner) url(id string) string {
	return path.Join(`/api/admin/config_repos`, url.PathEscape(id))
}

func (r *RmRunner) onSuccess(res *dub.Response) error {
	return api.ReadBodyAndDo(res, func(b []byte) error {
		if r, err := api.ParseMessage(b); err != nil {
			return utils.InspectError(err, `parsing config-repo delete response: %s`, string(b))
		} else {
			utils.Echofln(r.String())
		}
		return nil
	})
}

func (r *RmRunner) onFail(res *dub.Response) error {
	return api.ReadBodyAndDo(res, func(b []byte) error {
		_, id := path.Split(res.Raw.Request.URL.Path)
		api.DieOnAuthError(res)
		api.DieOnNotFound(res, `No such config-repo with id: %q`, id)

		return nil
	})
}

func init() {
	RootCmd.AddCommand(RmCmd)
}
