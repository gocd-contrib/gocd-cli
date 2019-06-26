package cmd

import (
	"github.com/gocd-contrib/gocd-cli/cfg"
	"github.com/gocd-contrib/gocd-cli/cmd/config"
	"github.com/gocd-contrib/gocd-cli/cmd/configrepo"
	"github.com/gocd-contrib/gocd-cli/cmd/encrypt"
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:       "gocd",
	Short:     "A command-line companion to a GoCD server",
	ValidArgs: []string{"config", "configrepo", "encrypt", "help"}, // bash-completion
}

var cfgFile string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		utils.AbortLoudly(err)
	}
}

func init() {
	cobra.OnInitialize(func() {
		if err := cfg.Setup(cfgFile); err != nil {
			utils.AbortLoudly(err)
		} else {
			utils.Debug("Loaded config from: %s", cfg.Conf().ConfigFile())
		}
	})

	RootCmd.AddCommand(config.RootCmd)
	RootCmd.AddCommand(configrepo.RootCmd)
	RootCmd.AddCommand(encrypt.RootCmd)
	RootCmd.AddCommand(AboutCommand)

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.gocd/settings.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&utils.SuppressOutput, "quiet", "q", false, "silence output")
	RootCmd.PersistentFlags().BoolVarP(&utils.DebugMode, "debug", "X", false, "debug output; overrides --quiet")
}
