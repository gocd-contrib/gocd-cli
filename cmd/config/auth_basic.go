package config

import (
	"github.com/spf13/cobra"
)

var BasicAuthCmd = &cobra.Command{
	Use:   "auth-basic <user> <pass>",
	Short: "Configures basic authentication for API requests to the GoCd server instance",
	Long:  "This sets the basic authentication credentials for GoCD API requests used by this CLI tool.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		runBasicAuth(args)
	},
}

func runBasicAuth(args []string) {
	// 1. get username and password arguments
	// 2. validate non-empty strings
	// 3. write auth to config (include type=basic + credentials)
	conf().SetBasicAuth("user", "pass")
}

func init() {
	RootCmd.AddCommand(BasicAuthCmd)
}
