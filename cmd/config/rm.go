package config

import (
	"strings"

	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var RmCmd = &cobra.Command{
	Use:   "rm <config-key>",
	Short: "Deletes a configured value",
	Example: strings.Trim(`
  gocd config rm auth          # Deletes the current authentication configuration
  gocd config rm server-url    # Deletes the server URL configuration`, "\n"),
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rmConfig.Run(args)
	},
}

var rmConfig = &RmConfigRunner{}

type RmConfigRunner struct{}

func (su *RmConfigRunner) Run(args []string) {
	if err := conf().Unset(args[0]); err != nil {
		utils.AbortLoudly(err)
	}
}

func init() {
	RootCmd.AddCommand(RmCmd)
}
