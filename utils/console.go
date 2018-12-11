package utils

import (
	"fmt"
	"os"
)

var SuppressOutput bool

// Writes to STDOUT unless SuppressOutput is set.
// Uses `printf()` formatting, but also appends a newline
// like `println()`.
func Echof(f string, t ...interface{}) {
	if !SuppressOutput {
		fmt.Fprintf(os.Stdout, f+"\n", t...)
	}
}

// Writes to STDERR unless SuppressOutput is set.
// Uses `printf()` formatting, but also appends a newline
// like `println()`.
func Errf(f string, t ...interface{}) {
	if !SuppressOutput {
		fmt.Fprintf(os.Stderr, f+"\n", t...)
	}
}

// Exits with exitCode after printing message.
// Automatically selects STDOUT vs STDERR depending
// on value of exitCode
func Die(exitCode int, f string, t ...interface{}) {
	if exitCode != 0 {
		Errf(f, t...)
	} else {
		Echof(f, t...)
	}

	os.Exit(exitCode)
}

// Prints `error.Error()` to STDERR and exits with failure
func Abort(err error) {
	Die(1, "%s\n", err)
}

// *Loudly counterparts will output messages whether SuppressOutput
// is set or not. Some messages (generally unexpected errors, or
// incorrect invocation) should always be outputted so as to avoid
// false negatives.

func DieLoudly(exitCode int, f string, t ...interface{}) {
	if exitCode != 0 {
		fmt.Fprintf(os.Stderr, f+"\n", t...)
	} else {
		fmt.Fprintf(os.Stdout, f+"\n", t...)
	}

	os.Exit(exitCode)

}

func AbortLoudly(err error) {
	DieLoudly(1, "%s\n", err)
}
