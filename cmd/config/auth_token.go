package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var TokenAuthCmd = &cobra.Command{
	Use:   "auth-token TOKEN",
	Short: "Configures token authentication for API requests to the GoCD server instance",
	Long:  "This sets the Personal Access Token for GoCD API requests made by this CLI tool.",
	Example: strings.Trim(`
  gocd config auth-token abcd1234             # sets the auth token to abcd1234
  cat token.txt | gocd config auth-token -    # sets the token from STDIN; must pass in the "-" argument`, "\n"),
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tokenAuth.Run(args)
	},
}

var tokenAuth = &TokenAuthRunner{}

type TokenAuthRunner struct{}

func (su *TokenAuthRunner) Run(args []string) {
	var token string
	if utils.HasShellPipe() {
		if `-` != args[0] {
			utils.DieLoudly(1, `For piped input, you must specify "-" as the argument`)
		}

		if b, err := ioutil.ReadAll(os.Stdin); err != nil {
			utils.DieLoudly(1, `Failed to read from STDIN; Cause: %v`, err)
		} else {
			token = strings.TrimSpace(string(b))
		}
	} else {
		token = args[0]
	}

	if err := conf().SetTokenAuth(token); err != nil {
		utils.AbortLoudly(err)
	}
}

func init() {
	RootCmd.AddCommand(TokenAuthCmd)
}
