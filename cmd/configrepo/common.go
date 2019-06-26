package configrepo

import (
	"github.com/gocd-contrib/gocd-cli/api"
	"github.com/gocd-contrib/gocd-cli/dub"
	"github.com/gocd-contrib/gocd-cli/utils"
)

func DieWhenPluginIdNotSet() {
	if "" == PluginId {
		utils.DieLoudly(1, "You must provide a --plugin-id")
	}
}

func printResponse(res *dub.Response) error {
	return api.ReadBodyAndDo(res, func(b []byte) error {
		utils.Echofln(string(b))
		return nil
	})
}
