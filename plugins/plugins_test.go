package plugins

import (
	"testing"
)

func TestPluginByIdReturnsErrorIfCantOpenFile(t *testing.T) {
	as := asserts(t)

	f, err := PluginById("test.config.plugin", "invalidpath")

	as.err("open invalidpath: no such file or directory", err)
	as.eq("", f)
}

func TestPluginById(t *testing.T) {
	as := asserts(t)

	found, err := PluginById("test.config.plugin", "testdata/testplugin.jar")
	as.ok(err)

	as.eq("testdata/testplugin.jar", found)
}

func TestPluginByIdWillSearchDirectory(t *testing.T) {
	as := asserts(t)

	found, err := PluginById("test.config.plugin", "testdata")
	as.ok(err)

	as.eq("testdata/testplugin.jar", found)
}

func TestPluginByIdReturnsErrorIfNoMatchingJar(t *testing.T) {
	as := asserts(t)

	found, err := PluginById("nomatch.id", "testdata")

	as.err(`No matching plugin jar with id "nomatch.id" found in path "testdata"`, err)
	_, matchesType := err.(*PluginNotFoundError)
	as.is(matchesType)
	as.eq("", found)
}

func TestPluginByIdReturnsErrorOnCorruptedPlugin(t *testing.T) {
	as := asserts(t)

	found, err := PluginById("whatever", "testdata/baddata/badplugin.jar")

	as.err("XML syntax error on line 2: unexpected EOF", err)
	as.eq("", found)
}
