package config

import (
	"strings"

	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete <config-key>",
	Short: "Deletes a configured value",
	Example: strings.Trim(`
  gocd config delete auth          # Deletes the current authentication configuration
  gocd config delete server-url    # Deletes the server URL configuration`, "\n"),
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteConfig.Run(args)
	},
}

var deleteConfig = &DeleteConfigRunner{}

type DeleteConfigRunner struct{}

func (su *DeleteConfigRunner) Run(args []string) {
	if err := conf().Unset(args[0]); err != nil {
		utils.AbortLoudly(err)
	}
}

func init() {
	RootCmd.AddCommand(DeleteCmd)
}
