package config

import (
	"strings"

	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var NoAuthCmd = &cobra.Command{
	Use:   "auth-none",
	Short: "Explicitly specifies that API requests require no authentication on the GoCD server instance",
	Long:  "Setting this ensures that no `Authorization` header is sent by this CLI tool; you MUST set this if security is disabled on your GoCD server.",
	Example: strings.Trim(`
  gocd config auth-none                       # explicitly specifies that API requests are not sent with credentials`, "\n"),
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		authNone.Run(args)
	},
}

var authNone = &AuthNoneRunner{}

type AuthNoneRunner struct{}

func (su *AuthNoneRunner) Run(args []string) {
	if err := conf().SetRequestsAreUnauthenticated(); err != nil {
		utils.AbortLoudly(err)
	}
}

func init() {
	RootCmd.AddCommand(NoAuthCmd)
}
