package configrepo

import (
	"github.com/gocd-contrib/gocd-cli/api"
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var RmCmd = &cobra.Command{
	Use:   "rm id",
	Short: "Deletes a config-repo by id",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rm.Run(args)
	},
}

var rm = &RmRunner{}

type RmRunner struct{}

func (r *RmRunner) Run(args []string) {
	if err := Model.DeleteRepo(args[0], func(msg api.MessageResponse) error {
		utils.Echofln(msg.String())
		return nil
	}); err != nil {
		utils.AbortLoudly(err)
	}
}

func init() {
	RootCmd.AddCommand(RmCmd)
}
