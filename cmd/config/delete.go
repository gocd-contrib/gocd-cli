package config

import (
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete <config-key>",
	Short: "Deletes a configured value",
	Args:  cobra.ExactArgs(1),
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
