package github

import (
	"fmt"
	"strings"

	"github.com/blang/semver"
	"github.com/gocd-contrib/gocd-cli/utils"
)

func ResolveVersionJar(rels []Release, filterBy string, stableOnly bool) (asset *Asset, err error) {
	if 0 == len(rels) {
		return nil, nil
	}

	if !stableOnly && "" == filterBy {
		return findJarAsset(&(rels[0]))
	} else {
		filt := buildFilter(filterBy, !stableOnly)
		for _, r := range rels {
			if filt(&r) {
				return findJarAsset(&r)
			}
		}
	}

	return nil, fmt.Errorf("Cannot find a release that matches version spec `%s` and stable-only=%t", filterBy, stableOnly)
}

func buildFilter(filterSpec string, allowPrerelease bool) func(*Release) bool {
	if r, err := semver.ParseRange(filterSpec); err != nil {
		utils.DieLoudly(1, "Don't know how to parse version spec `%s`: %v", filterSpec, err)
		return nil
	} else {
		return func(rel *Release) bool {
			if v, err := semver.Parse(rel.Version); err != nil {
				utils.DieLoudly(1, "Cannot parse release version `%s`; does not conform to semantic version (%v)", err)
				return false
			} else {
				return r(v) && (allowPrerelease || !rel.Prerelease)
			}
		}
	}
}

func findJarAsset(rel *Release) (*Asset, error) {
	for _, j := range rel.Assets {
		if strings.HasSuffix(j.Name, ".jar") {
			return &j, nil
		}
	}

	return nil, fmt.Errorf("Could not resolve jar asset for version %s", rel.Version)
}
