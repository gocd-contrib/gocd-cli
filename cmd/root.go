package cmd

import (
	"github.com/gocd-contrib/gocd-cli/cmd/config"
	"github.com/gocd-contrib/gocd-cli/cmd/configrepo"
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:       "gocd",
	Short:     "A command-line companion to a GoCD server",
	Long:      `A command-line helper to GoCD to help build config-repos, among other things (?)`,
	ValidArgs: []string{"config", "configrepo", "help"}, // bash-completion
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.AbortLoudly(err)
	}
}

func init() {
	cobra.OnInitialize(config.Setup)
	rootCmd.AddCommand(config.RootCmd)
	rootCmd.AddCommand(configrepo.RootCmd)

	rootCmd.PersistentFlags().StringVarP(&config.CfgFile, "config", "c", "", "config file (default is $HOME/.gocd/settings.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&utils.SuppressOutput, "quiet", "q", false, "silence output")
}
