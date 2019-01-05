package config

import (
	"os"
	"path/filepath"

	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	CONFIG_DIRNAME  = ".gocd"
	CONFIG_FILENAME = "settings"
	CONFIG_ENV_PFX  = "gocdcli"
)

var CfgFile string
var Config = viper.New()

var RootCmd = &cobra.Command{
	Use:       "config",
	Short:     "GoCD CLI configuration",
	ValidArgs: []string{"auth-basic", "server-url", "help", "clear"}, // bash-completion
}

func Setup() {
	Config.SetEnvPrefix(CONFIG_ENV_PFX)
	Config.AutomaticEnv() // read in environment variables that match

	if CfgFile != "" {
		// Use config file from the flag.
		Config.SetConfigFile(CfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			utils.AbortLoudly(err)
		}

		cfgDir := filepath.Join(home, CONFIG_DIRNAME)
		os.MkdirAll(cfgDir, os.ModePerm)

		Config.AddConfigPath(cfgDir)
		Config.SetConfigName(CONFIG_FILENAME)
	}

	// If a config file is found, read it in.
	if err := Config.ReadInConfig(); err == nil {
		utils.Echofln("Using config file:", Config.ConfigFileUsed())
	}
}
