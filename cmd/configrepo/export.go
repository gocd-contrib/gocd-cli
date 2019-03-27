package configrepo

import (
	"io/ioutil"
	"mime"
	"net/url"

	"github.com/gocd-contrib/gocd-cli/api"
	"github.com/gocd-contrib/gocd-cli/dub"
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var ExportCmd = &cobra.Command{
	Use:     "export <pipeline name>",
	Aliases: []string{"ex"},
	Short:   "Exports the specified pipeline as a config-repo definition in the indicated config-repo plugin format",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		export.Run(args)
	},
}

var export = &ExportRunner{}

type ExportRunner struct {
	UseStdout bool
}

func (er *ExportRunner) Run(args []string) {
	if "" == PluginId {
		utils.DieLoudly(1, "You must provide a --plugin-id")
	}

	if err := api.V1.Get(er.url(args[0])).Send(er.onSuccess, er.onFail); err != nil {
		utils.AbortLoudly(err)
	}
}

func (er *ExportRunner) url(pipeline string) string {
	return dub.AddQuery(`/api/admin/export/pipelines/`+url.PathEscape(pipeline), url.Values{
		`plugin_id`: {PluginId},
	})
}

func (er *ExportRunner) onSuccess(res *dub.Response) error {
	return api.ReadBodyAndDo(res, func(data []byte) error {
		if er.UseStdout {
			utils.Echofln(string(data))
			return nil
		}

		if _, params, err := mime.ParseMediaType(res.Headers.Get(`Content-Disposition`)); err == nil {
			return ioutil.WriteFile(params[`filename`], data, 0644)
		} else {
			return err
		}
	})
}

func (er *ExportRunner) onFail(res *dub.Response) error {
	return api.ReadBodyAndDo(res, func(data []byte) error {
		if m, err := api.ParseMessage(data); err == nil {
			utils.Die(1, m.String()) // call Die() here instead of letting it hit AbortLoudly so it can be silenced
			return nil
		} else {
			return err
		}
	})
}

func init() {
	RootCmd.AddCommand(ExportCmd)
	ExportCmd.Flags().BoolVar(&export.UseStdout, "stdout", false, "print to STDOUT instead of a file")
}
