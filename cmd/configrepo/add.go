package configrepo

import (
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Creates a new config-repo",
}

func init() {
	AddCmd.AddCommand(GitRepoCmd(true))
	AddCmd.AddCommand(HgRepoCmd(true))
	AddCmd.AddCommand(SvnRepoCmd(true))
	AddCmd.AddCommand(P4RepoCmd(true))
	AddCmd.AddCommand(TfsRepoCmd(true))

	RootCmd.AddCommand(AddCmd)
}
