package config

import (
	"github.com/gocd-contrib/gocd-cli/cfg"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:       "config",
	Aliases:   []string{"cf"},
	Short:     "GoCD CLI configuration",
	ValidArgs: []string{"auth-basic", "server-url", "help", "delete"}, // bash-completion
}

// convenvience method so subcommands don't need to import cfg
func conf() *cfg.Config {
	return cfg.Conf()
}
