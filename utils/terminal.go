package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func StdoutOrDevNull() io.Writer {
	if !DebugMode && SuppressOutput {
		return ioutil.Discard
	}
	return os.Stdout
}

func StderrOrDevNull() io.Writer {
	if !DebugMode && SuppressOutput {
		return ioutil.Discard
	}
	return os.Stderr
}

// Detects if this invocation has piped STDIN
func HasShellPipe() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}

func UseXargsOverPipe(rawArgs []string) error {
	if HasShellPipe() {
		return &MustUseXargs{Invocation: rawArgs}
	}
	return nil
}

type MustUseXargs struct {
	Invocation []string
}

func (e *MustUseXargs) Error() string {
	command := strings.Join(e.Invocation, ` `)

	f := strings.Join([]string{
		"This command does not directly accept a shell pipe; perhaps you meant to use `xargs`?",
		``,
		`Example (*Nix):`,
		"  `<piped-output> | xargs %s`",
		``,
		`Example (PowerShell):`,
		"  `<piped-output> | %% { $_.split()[0] } | %% { %s $_ }`",
	}, "\n") + "\n"

	return fmt.Sprintf(f, command, command)
}
