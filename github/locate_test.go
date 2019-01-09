package github

import (
	"testing"
)

func jarname(version string) string {
	return "test.name-" + version + ".jar"
}

func jarurl(version string) string {
	return "http://test.com/releases/download/" + version + "/" + jarname(version)
}

func createRelease(version string, prerelease bool) Release {
	return Release{Version: version, Prerelease: prerelease, Assets: []Asset{Asset{Name: jarname(version), Url: jarurl(version)}}}
}

func TestResolveVersionJar(t *testing.T) {
	as := asserts(t)
	releases := []Release{createRelease("1.1.1", true)}

	asset, err := ResolveVersionJar(releases, "1.1.1", false)
	as.ok(err)

	as.eq(Asset{Name: jarname("1.1.1"), Url: jarurl("1.1.1")}, *asset)

	asset, err = ResolveVersionJar(releases, "", false)
	as.ok(err)

	as.eq(Asset{Name: jarname("1.1.1"), Url: jarurl("1.1.1")}, *asset)
}

func TestResolveVersionJarCanSkipPreReleaseJars(t *testing.T) {
	as := asserts(t)
	releases := []Release{createRelease("1.1.0", false), createRelease("1.1.1", true)}

	asset, err := ResolveVersionJar(releases, "1.1.x", true)
	as.ok(err)

	as.eq(jarname("1.1.0"), asset.Name)
}

func TestResolveVersionJarCanFindVersionWithinRange(t *testing.T) {
	as := asserts(t)
	version := "0.8.2"
	releases := []Release{createRelease(version, false)}

	asset, err := ResolveVersionJar(releases, ">=0.5.0 <0.8.0 || >=0.8.1 !0.8.3", false)
	as.ok(err)

	as.eq(Asset{Name: jarname(version), Url: jarurl(version)}, *asset)
}

func TestResolveVersionJarReturnsErrorIfNoMatchingJar(t *testing.T) {
	as := asserts(t)
	versionFilter := ">=0.5.0 <0.8.0"
	releases := []Release{createRelease("0.8.2", false)}

	_, err := ResolveVersionJar(releases, versionFilter, false)
	as.err("Cannot find a release that matches version spec `"+versionFilter+"` and stable-only=false", err)
}

func TestResolveVersionJarReturnsErrorIfCannotParseFilterVersion(t *testing.T) {
	as := asserts(t)
	_, err := ResolveVersionJar([]Release{createRelease("1.1.1", true)}, "mumbojumborelease", true)
	as.err("Don't know how to parse version spec `mumbojumborelease`: Could not get version from string: \"mumbojumborelease\"", err)
}

func TestResolveVersionJarReturnsErrorIfCannotParseReleaseVersion(t *testing.T) {
	as := asserts(t)
	_, err := ResolveVersionJar([]Release{createRelease("mumbojumborelease", true)}, "1.1.1", true)
	as.err("Cannot parse release version `mumbojumborelease`; does not conform to semantic version (No Major.Minor.Patch elements found)", err)
}

func TestResolveVersionJarReturnErrorIfCannotFindAssetJarForRelease(t *testing.T) {
	as := asserts(t)
	release := createRelease("1.1.1", false)
	release.Assets[0].Name = "not a jar"

	_, err := ResolveVersionJar([]Release{release}, "1.1.1", true)
	as.err("Could not resolve jar asset for version 1.1.1", err)
}

func TestResolveVersionJarReturnErrorIfThereAreNoReleases(t *testing.T) {
	as := asserts(t)
	_, err := ResolveVersionJar([]Release{}, "", true)
	as.err("Could not find any releases to filter", err)
}
