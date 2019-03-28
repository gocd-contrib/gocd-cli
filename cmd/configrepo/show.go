package configrepo

import (
	"net/url"
	"path"

	"github.com/gocd-contrib/gocd-cli/api"
	"github.com/gocd-contrib/gocd-cli/dub"
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var ShowCmd = &cobra.Command{
	Use:   "show id",
	Short: "Displays the settings for an existing config-repo",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		show.Run(args)
	},
}

var show = &ShowRunner{}

type ShowRunner struct{}

func (r *ShowRunner) Run(args []string) {
	if err := api.V1.Get(r.url(args[0])).Send(r.onSuccess, r.onFail); err != nil {
		utils.AbortLoudly(err)
	}
}

func (r *ShowRunner) url(id string) string {
	return path.Join(`/api/admin/config_repos`, url.PathEscape(id))
}

func (r *ShowRunner) onSuccess(res *dub.Response) error {
	return api.ReadBodyAndDo(res, func(b []byte) error {
		utils.Echofln(string(b))
		return nil
	})
}

func (r *ShowRunner) onFail(res *dub.Response) error {
	return api.ReadBodyAndDo(res, func(b []byte) error {
		_, id := path.Split(res.Raw.Request.URL.Path)
		api.DieOnAuthError(res)
		api.DieOnNotFound(res, `No such config-repo with id: %q`, id)

		return nil
	})
}

func init() {
	// ShowCmd.Flags().BoolVar(&show.Raw, "raw", false, "machine-readable output (JSON)")
	RootCmd.AddCommand(ShowCmd)
}
