package configrepo

import (
	"log"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var PluginId string
var PluginDir string
var PluginJar string

// RootCmd represents the configrepo command
var RootCmd = &cobra.Command{
	Use:   "configrepo",
	Short: "GoCD config-repo functions",
	Long:  `Functions to help development of config-repos in GoCD (pipeline configs as code)`,
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&PluginDir, "plugin-dir", "d", "", "The plugin directory to search for plugins")

	RootCmd.PersistentFlags().StringVarP(&PluginId, "plugin-id", "i", "", "The config-repo plugin to use (e.g., yaml.config.plugin)")
	RootCmd.MarkFlagRequired("plugin-id")

	if PluginDir == "" {
		if d, err := homedir.Dir(); err == nil {
			PluginDir = filepath.Join(d, ".gocd", "plugins")
		} else {
			log.Fatal(err)
		}
	}

	if err := os.MkdirAll(PluginDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}
