package config

import (
	"github.com/spf13/cobra"
)

var ServerUrlCmd = &cobra.Command{
	Use:   "server-url",
	Short: "Configures the base url for API requests to the GoCD server instance",
	Long:  "This sets the base url for GoCD API requests used by this CLI tool. The base url includes the protocol, host, port (if applicable), and path (anything before /go, if applicable).",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runServerUrl(args)
	},
}

func runServerUrl(args []string) {
	// 1. get the base url string from args; probably need to validate that it's not empty
	// 2. want to parse/validate/normalize the URL using golang's "net/url" package
	// 3. write parsed url to config
	//      - either as a full url
	//      - or as individual parts (e.g. host, port, protocol, path, authString)
}

func init() {
	RootCmd.AddCommand(ServerUrlCmd)
}
