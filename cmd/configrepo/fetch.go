package configrepo

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocd-contrib/gocd-cli/dub"
	"github.com/gocd-contrib/gocd-cli/github"
	"github.com/gocd-contrib/gocd-cli/plugins"
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var FetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetches configrepo plugins",
	Run: func(cmd *cobra.Command, args []string) {
		runFetch(args)
	},
}

var stableOnly bool
var filterBy string

func runFetch(args []string) {
	if "" == PluginId {
		utils.DieLoudly(1, "You must provide a --plugin-id")
	}

	if _, err := FetchPlugin(PluginId); err != nil {
		utils.AbortLoudly(err)
	}
}

func FetchPlugin(id string) (string, error) {
	releases := make([]github.Release, 0)

	if err := dub.New().Get(url(PluginId)).Do(func(res *dub.Response) error {
		payload, err := res.ReadAll()

		if err != nil {
			return err
		}

		return json.Unmarshal(payload, &releases)
	}); nil != err {
		return "", err
	}

	if 0 == len(releases) {
		return "", fmt.Errorf("There are no available releases for %s", id)
	}

	if a, err := github.ResolveVersionJar(releases, filterBy, stableOnly); err != nil {
		return "", err
	} else {
		if existing := plugins.PluginById(PluginId, PluginDir); "" != existing {
			if utils.IsDir(existing) {
				utils.Errfln("[WARNING] `%s` is a directory; will not remove this, but please inspect.", existing)
			}

			if utils.IsFile(existing) {
				utils.Echofln("Removing existing %s plugin %s", PluginId, existing)
				os.RemoveAll(existing)
			}
		}

		return utils.Wget(a.Url, a.Name, PluginDir)
	}
}

func url(pluginId string) (releaseUrl string) {
	if v, ok := plugins.ConfigRepo[PluginId]; ok {
		releaseUrl = v.Url
	} else {
		utils.DieLoudly(1, "Don't know how to fetch plugin `%s`; known plugins: %s", pluginId, plugins.ConfigRepo.ShortList())
	}
	return
}

func init() {
	FetchCmd.Flags().BoolVar(&stableOnly, "stable", false, "Restrict to stable (i.e., non-prerelease) releases")
	FetchCmd.Flags().StringVar(&filterBy, "match-version", "", "Specify a semver exact match, range (e.g., >=1.0.0 <2.0.0 || >=3.0.0 !3.0.1-beta.1), or wildcard (e.g., 0.8.x)")
	RootCmd.AddCommand(FetchCmd)
}
