package utils

import (
	"fmt"
	"os"
)

var SuppressOutput bool
var DebugMode bool

// Chainable/pass-thru error debug logger
func InspectError(err error, label string, t ...interface{}) error {
	if err == nil {
		return nil
	}

	if len(t) > 0 {
		label = fmt.Sprintf(label, t...)
	}

	Debug("While [%s], caught error:\n  -> %q", label, err.Error())

	return err
}

func Debug(f string, t ...interface{}) {
	if DebugMode {
		if len(t) > 0 {
			fmt.Printf(`[DEBUG] `+f+"\n", t...)
		} else {
			fmt.Println(`[DEBUG] ` + f)
		}
	}
}

// Writes to STDOUT unless SuppressOutput is set.
// Uses `printf()` formatting
func Echof(f string, t ...interface{}) {
	if DebugMode || !SuppressOutput {
		if len(t) > 0 {
			fmt.Fprintf(os.Stdout, f, t...)
		} else {
			fmt.Fprint(os.Stdout, f)
		}
	}
}

func Echofln(f string, t ...interface{}) {
	Echof(f+"\n", t...)
}

// Writes to STDERR unless SuppressOutput is set.
// Uses `printf()` formatting
func Errf(f string, t ...interface{}) {
	if DebugMode || !SuppressOutput {
		if len(t) > 0 {
			fmt.Fprintf(os.Stderr, f, t...)
		} else {
			fmt.Fprint(os.Stderr, f)
		}
	}
}

func Errfln(f string, t ...interface{}) {
	Errf(f+"\n", t...)
}

// Exits with exitCode after printing message.
// Automatically selects STDOUT vs STDERR depending
// on value of exitCode
func Die(exitCode int, f string, t ...interface{}) {
	if exitCode != 0 {
		Errfln(f, t...)
	} else {
		Echofln(f, t...)
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

func AbortLoudlyOnError(err error) {
	if err != nil {
		AbortLoudly(err)
	}
}

func AbortLoudly(err error) {
	DieLoudly(1, "%s", err.Error())
}
