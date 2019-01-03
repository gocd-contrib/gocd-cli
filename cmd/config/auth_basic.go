package config

import (
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var BasicAuthCmd = &cobra.Command{
	Use:   "auth-basic <user> <pass>",
	Short: "Configures basic authentication for API requests to the GoCd server instance",
	Long:  "This sets the basic authentication credentials for GoCD API requests used by this CLI tool.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := runBasicAuth(args); err != nil {
			utils.AbortLoudly(err)
		}
	},
}

func runBasicAuth(args []string) error {
	return conf().SetBasicAuth(args[0], args[1])
}

func init() {
	RootCmd.AddCommand(BasicAuthCmd)
}
