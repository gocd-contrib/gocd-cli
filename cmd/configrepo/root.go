package configrepo

import (
	"os"
	"path/filepath"

	"github.com/gocd-contrib/gocd-cli/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var PluginId string
var PluginDir string
var PluginJar string

// RootCmd represents the configrepo command
var RootCmd = &cobra.Command{
	Use:       "configrepo",
	Aliases:   []string{"cr"},
	Short:     "GoCD config-repo functions",
	Long:      `Functions to help development of config-repos in GoCD (pipeline configs as code)`,
	ValidArgs: []string{"syntax", "fetch", "preflight", "help"}, // bash-completion
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&PluginDir, "plugin-dir", "d", "", "The plugin directory to search for plugins")

	RootCmd.PersistentFlags().StringVarP(&PluginId, "plugin-id", "i", "", "The config-repo plugin to use (e.g., yaml.config.plugin)")
	RootCmd.MarkFlagRequired("plugin-id")

	// Alias flags for --plugin-id
	RootCmd.PersistentFlags().VarPF(newJsonFlag(false), "json", "", "Alias for '--plugin-id json.config.plugin'").NoOptDefVal = `true`
	RootCmd.PersistentFlags().VarPF(newYamlFlag(false), "yaml", "", "Alias for '--plugin-id yaml.config.plugin'").NoOptDefVal = `true`
	RootCmd.PersistentFlags().VarPF(newGroovyFlag(false), "groovy", "", "Alias for '--plugin-id cd.go.contrib.plugins.configrepo.groovy'").NoOptDefVal = `true`

	if PluginDir == "" {
		if d, err := homedir.Dir(); err == nil {
			PluginDir = filepath.Join(d, ".gocd", "plugins")
		} else {
			utils.AbortLoudly(err)
		}
	}

	if err := os.MkdirAll(PluginDir, os.ModePerm); err != nil {
		utils.AbortLoudly(err)
	}
}
