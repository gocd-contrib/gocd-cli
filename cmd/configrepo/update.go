package configrepo

import (
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Modifies an existing config-repo",
}

func init() {
	UpdateCmd.AddCommand(GitRepoCmd(false))
	UpdateCmd.AddCommand(HgRepoCmd(false))
	UpdateCmd.AddCommand(SvnRepoCmd(false))
	UpdateCmd.AddCommand(P4RepoCmd(false))
	UpdateCmd.AddCommand(TfsRepoCmd(false))

	RootCmd.AddCommand(UpdateCmd)
}
