package plugins

import (
	"fmt"
	"strings"

	"github.com/blang/semver"
	"github.com/gocd-contrib/gocd-cli/utils"
)

type PluginMap map[string]*Info

func (pm PluginMap) Ids() []string {
	keys := make([]string, len(pm))
	i := 0
	for k := range pm {
		keys[i] = k
		i++
	}
	return keys
}

func (pm PluginMap) ShortList() string {
	return fmt.Sprintf("[%s]", strings.Join(pm.Ids(), ", "))
}

type Info struct {
	Url     string
	Version string
	Compat  semver.Range
}

func (info *Info) IsCompatible(version string) bool {
	if v, err := semver.Parse(version); err == nil {
		return info.Compat(v)
	} else {
		utils.AbortLoudly(err)
		return false
	}
}

func NewInfo(url string, version string) *Info {
	sv, err := semver.ParseRange(version)
	if err != nil {
		utils.AbortLoudly(err)
	}
	return &Info{Url: url, Version: version, Compat: sv}
}
