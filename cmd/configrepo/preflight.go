package configrepo

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"

	"github.com/gocd-contrib/gocd-cli/api"
	"github.com/gocd-contrib/gocd-cli/dub"
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var PreflightCmd = &cobra.Command{
	Use:   "preflight",
	Short: "Preflights any number of definition files for syntax, structure, and dependencies against a running GoCD server",
	Run: func(cmd *cobra.Command, args []string) {
		preflight.Run(args)
	},
}

var preflight = &PreflightRunner{}

type PreflightRunner struct {
	RepoId string
}

func (pr *PreflightRunner) Run(args []string) {
	if "" == PluginId {
		utils.DieLoudly(1, "You must provide a --plugin-id")
	}

	body := dub.NewPipedMultipart()

	for _, f := range args {
		body.AddFile(`files[]`, f)
	}

	if err := api.V1.Post(pr.url(), body).Send(pr.handleApiResponse); err != nil {
		utils.AbortLoudly(err)
	}
}

func (pr *PreflightRunner) handleApiResponse(res *dub.Response) error {
	res.Consume(func(reader io.Reader) error {
		if b, err := ioutil.ReadAll(reader); err != nil {
			return utils.InspectError(err, `reading preflight response from %q`, res.Raw.Request.URL)
		} else {
			if res.IsAuthError() {
				utils.DieLoudly(1, `Invalid credentials. Either the username or password configured is incorrect`)
			}

			if res.IsError() {
				if msg, err := api.ParseMessage(b); err == nil {
					return fmt.Errorf(`Unexpected response %d: %s`, res.Status, msg)
				} else {
					return utils.InspectError(err, `parsing api error %d response: %q`, res.Status, string(b))
				}
			}

			if result, err := ParseCrPreflight(b); err == nil {
				if result.Valid {
					utils.Echofln(`OK`)
				} else {
					utils.Die(1, result.DisplayErrors())
				}
			} else {
				return utils.InspectError(err, `parsing preflight api response %q`, string(b))
			}
		}

		return nil
	})

	return nil
}

func (pr *PreflightRunner) url() string {
	query := url.Values{
		`pluginId`: {PluginId},
	}

	if "" != pr.RepoId {
		query.Add(`repoId`, pr.RepoId)
	}

	return dub.AddQuery(`/api/admin/config_repo_ops/preflight`, query)
}

func ParseCrPreflight(body []byte) (*api.CrPreflightResponse, error) {
	r := &api.CrPreflightResponse{}
	if err := json.Unmarshal(body, r); err == nil {
		return r, nil
	} else {
		return nil, err
	}
}

func init() {
	RootCmd.AddCommand(PreflightCmd)
	PreflightCmd.Flags().StringVarP(&preflight.RepoId, "repo-id", "r", "", "A config-repo ID; use this preflighting change to an existing config-repo")
}
