package utils

import (
	"io"
	"io/ioutil"
	"os"
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
