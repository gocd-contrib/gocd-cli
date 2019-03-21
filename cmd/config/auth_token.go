package config

import (
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var TokenAuthCmd = &cobra.Command{
	Use:   "auth-token <token>",
	Short: "Configures token authentication for API requests to the GoCD server instance",
	Long:  "This sets Personal Access Token for GoCD API requests used by this CLI tool.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tokenAuth.Run(args)
	},
}

var tokenAuth = &TokenAuthRunner{}

type TokenAuthRunner struct{}

func (su *TokenAuthRunner) Run(args []string) {
	if err := conf().SetTokenAuth(args[0]); err != nil {
		utils.AbortLoudly(err)
	}
}

func init() {
	RootCmd.AddCommand(TokenAuthCmd)
}
