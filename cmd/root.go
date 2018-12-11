package cmd

import (
	"os"
	"path/filepath"

	"github.com/gocd-contrib/gocd-cli/cmd/configrepo"
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "gocd",
	Short: "A command-line companion to a GoCD server",
	Long:  `A command-line helper to GoCD to help build config-repos, among other things (?)`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.AbortLoudly(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(configrepo.RootCmd)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.gocd/settings.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&utils.SuppressOutput, "quiet", "q", false, "silence output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			utils.AbortLoudly(err)
		}

		cfgDir := filepath.Join(home, ".gocd")
		os.MkdirAll(cfgDir, os.ModePerm)

		viper.AddConfigPath(cfgDir)
		viper.SetConfigName("settings")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		utils.Echof("Using config file:", viper.ConfigFileUsed())
	}
}
