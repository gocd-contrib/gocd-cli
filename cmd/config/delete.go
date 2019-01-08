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
		if err := runDelete(args); err != nil {
			utils.AbortLoudly(err)
		}
	},
}

func runDelete(args []string) error {
	return conf().Unset(args[0])
}

func init() {
	RootCmd.AddCommand(DeleteCmd)
}
