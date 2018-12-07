package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/marques-work/gocd-cli/cmd/configrepo"
	homedir "github.com/mitchellh/go-homedir"
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
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(configrepo.RootCmd)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.gocd/settings.yaml)")
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
			fmt.Println(err)
			os.Exit(1)
		}

		cfgDir := filepath.Join(home, ".gocd")
		os.MkdirAll(cfgDir, os.ModePerm)

		// Search config in home directory with name ".clicli" (without extension).
		viper.AddConfigPath(cfgDir)
		viper.SetConfigName("settings")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
