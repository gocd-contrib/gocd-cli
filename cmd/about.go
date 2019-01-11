package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gocd-contrib/gocd-cli/meta"
	lol "github.com/kris-nova/lolgopher"
	"github.com/spf13/cobra"
)

const (
	BANNER = `
A command-line companion to:

_____/\\\\\\\\\\\\______________________/\\\\\\\\\__/\\\\\\\\\\\\____
 ___/\\\//////////____________________/\\\////////__\/\\\////////\\\__
  __/\\\_____________________________/\\\/___________\/\\\______\//\\\_
   _\/\\\____/\\\\\\\_____/\\\\\_____/\\\_____________\/\\\_______\/\\\_
    _\/\\\___\/////\\\___/\\\///\\\__\/\\\_____________\/\\\_______\/\\\_
     _\/\\\_______\/\\\__/\\\__\//\\\_\//\\\____________\/\\\_______\/\\\_
      _\/\\\_______\/\\\_\//\\\__/\\\___\///\\\__________\/\\\_______/\\\__
       _\//\\\\\\\\\\\\/___\///\\\\\/______\////\\\\\\\\\_\/\\\\\\\\\\\\/___
        __\////////////_______\/////___________\/////////__\////////////_____


           gocd-cli v%s

`
)

var AboutCommand = &cobra.Command{
	Use:   "about",
	Short: "About GoCD CLI",
	Run: func(cmd *cobra.Command, args []string) {
		ct := strings.ToLower(os.Getenv(`COLORTERM`))
		var w io.Writer

		if `truecolor` == ct || `24bit` == ct {
			w = lol.NewTruecolorLolWriter()
		} else {
			w = lol.NewLolWriter()
		}

		fmt.Fprintf(w, BANNER, meta.VersionString())
	},
}
