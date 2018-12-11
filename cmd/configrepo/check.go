package configrepo

import (
	"os"
	"os/exec"

	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks a definition file for syntactical and structural correctness",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runCheck(args)
	},
}

func runCheck(args []string) {
	if "" == PluginId {
		utils.DieLoudly(1, "You must provide a --plugin-id")
	}

	PluginJar = utils.LocatePlugin(PluginId, PluginDir)

	cmdArgs := append([]string{"-jar", PluginJar, "syntax"}, args...)
	utils.Echof("args: %v", cmdArgs)
	cmd := exec.Command("java", cmdArgs...)

	if !utils.ExecQ(cmd) {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(CheckCmd)
}
