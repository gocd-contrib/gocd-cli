package configrepo

import (
	"os"
	"os/exec"

	"github.com/marques-work/gocd-cli/utils"
	"github.com/spf13/cobra"
)

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks a definition file for syntactical and structural correctness",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		run(args)
	},
}

func run(args []string) {
	PluginJar = utils.LocatePlugin(PluginId, PluginDir)

	cmdArgs := append([]string{"-jar", PluginJar, "syntax"}, args...)
	cmd := exec.Command("java", cmdArgs...)

	if !utils.Exec(cmd, os.Stdout, os.Stderr) {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(CheckCmd)
}
