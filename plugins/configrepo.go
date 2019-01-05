package plugins

var ConfigRepo = PluginMap{
	"yaml.config.plugin": NewInfo("https://api.github.com/repos/tomzo/gocd-yaml-config-plugin/releases", ">=0.8.3"),
	"json.config.plugin": NewInfo("https://api.github.com/repos/tomzo/gocd-json-config-plugin/releases", ">=0.3.3"),
}
