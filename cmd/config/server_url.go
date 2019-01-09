package config

import (
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var ServerUrlCmd = &cobra.Command{
	Use:   "server-url",
	Short: "Configures the base url for API requests to the GoCD server instance",
	Long:  "This sets the base url for GoCD API requests used by this CLI tool. The base url includes the protocol, host, port (if applicable), and path (anything before /go, if applicable).",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverUrl.Run(args)
	},
}

var serverUrl = &ServerUrlRunner{}

type ServerUrlRunner struct{}

func (su *ServerUrlRunner) Run(args []string) {
	if err := conf().SetServerUrl(args[0]); err != nil {
		utils.AbortLoudly(err)
	}
}

func init() {
	RootCmd.AddCommand(ServerUrlCmd)
}
