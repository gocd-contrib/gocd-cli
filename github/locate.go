package github

import (
	"fmt"
	"strings"

	"github.com/blang/semver"
)

func ResolveVersionJar(rels []Release, filterBy string, stableOnly bool) (*Asset, error) {
	if 0 == len(rels) {
		return nil, fmt.Errorf("Could not find any releases to filter")
	}

	if !stableOnly && "" == filterBy {
		return findJarAsset(&(rels[0]))
	} else {
		if filt, err := buildFilter(filterBy, !stableOnly); err != nil {
			return nil, err
		} else {
			for _, r := range rels {
				ok, e2 := filt(&r)

				if e2 != nil {
					return nil, e2
				}

				if ok {
					return findJarAsset(&r)
				}
			}
		}
	}

	return nil, fmt.Errorf("Cannot find a release that matches version spec `%s` and stable-only=%t", filterBy, stableOnly)
}

func buildFilter(filterSpec string, allowPrerelease bool) (func(*Release) (bool, error), error) {
	if r, err := semver.ParseRange(filterSpec); err != nil {
		return nil, fmt.Errorf("Don't know how to parse version spec `%s`: %v", filterSpec, err)
	} else {
		return func(rel *Release) (bool, error) {
			if v, err := semver.Parse(rel.Version); err != nil {
				return false, fmt.Errorf("Cannot parse release version `%s`; does not conform to semantic version (%v)", rel.Version, err)
			} else {
				return r(v) && (allowPrerelease || !rel.Prerelease), nil
			}
		}, nil
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
