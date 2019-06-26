package config

import (
	"github.com/gocd-contrib/gocd-cli/cfg"
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:       "config",
	Aliases:   []string{"cf"},
	Short:     "GoCD CLI configuration",
	ValidArgs: []string{"auth-basic", "server-url", "help", "rm"}, // bash-completion
	Run: func(cmd *cobra.Command, args []string) {
		if !root.interactive {
			if err := cmd.Usage(); err != nil {
				utils.AbortLoudly(err)
			}
		} else {
			root.Run()
		}
	},
}

var root = &rootRunner{}

type rootRunner struct {
	interactive bool
}

func (r *rootRunner) Run() {

}

func init() {
	RootCmd.Flags().BoolVarP(&root.interactive, `interactive`, `i`, false, `Interactively configure settings`)
}

// convenvience method so subcommands don't need to import cfg
func conf() *cfg.Config {
	return cfg.Conf()
}
