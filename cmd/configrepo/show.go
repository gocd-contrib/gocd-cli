package configrepo

import (
	"encoding/json"

	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var ShowCmd = &cobra.Command{
	Use:   "show id",
	Short: "Displays the settings for an existing config-repo",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		show.Run(args)
	},
}

var show = &ShowRunner{}

type ShowRunner struct{}

func (r *ShowRunner) Run(args []string) {
	if err := Model.FetchRepo(args[0], func(repo *ConfigRepo) error {
		if b, err := json.MarshalIndent(repo, ``, `  `); err == nil {
			utils.Echofln(string(b))
			return nil
		} else {
			return err
		}
	}); err != nil {
		utils.AbortLoudly(err)
	}
	return
}

func init() {
	// ShowCmd.Flags().BoolVar(&show.Raw, "raw", false, "machine-readable output (JSON)")
	RootCmd.AddCommand(ShowCmd)
}
