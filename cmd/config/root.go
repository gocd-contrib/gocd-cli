package config

import (
	"github.com/gocd-contrib/gocd-cli/cfg"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:       "config",
	Aliases:   []string{"cf"},
	Short:     "GoCD CLI configuration",
	ValidArgs: []string{"auth-token", "auth-basic", "auth-none", "server-url", "help", "rm"}, // bash-completion
}

// convenvience method so subcommands don't need to import cfg
func conf() *cfg.Config {
	return cfg.Conf()
}
