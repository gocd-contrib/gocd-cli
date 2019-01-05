package configrepo

import (
	"os"
	"os/exec"

	"github.com/gocd-contrib/gocd-cli/plugins"
	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var SyntaxCmd = &cobra.Command{
	Use:   "syntax",
	Short: "Checks one or more definition files for syntactical correctness",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runCheck(args)
	},
}

func runCheck(args []string) {
	if "" == PluginId {
		utils.DieLoudly(1, "You must provide a --plugin-id")
	}

	PluginJar = plugins.LocatePlugin(PluginId, PluginDir)

	cmdArgs := append([]string{"-jar", PluginJar, "syntax"}, args...)
	cmd := exec.Command("java", cmdArgs...)

	if !utils.ExecQ(cmd) {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(SyntaxCmd)
}
