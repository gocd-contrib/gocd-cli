package utils

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func Stdout() io.Writer {
	if !DebugMode && SuppressOutput {
		return ioutil.Discard
	}
	return os.Stdout
}

func Stderr() io.Writer {
	if !DebugMode && SuppressOutput {
		return ioutil.Discard
	}
	return os.Stderr
}

func ExecQ(cmd *exec.Cmd) bool {
	return Exec(cmd, os.Stdin, Stdout(), Stderr())
}

func Exec(cmd *exec.Cmd, pipeIn io.Reader, pipeOut io.Writer, pipeErr io.Writer) bool {
	cmd.Stdin = pipeIn
	cmd.Stdout = pipeOut
	cmd.Stderr = pipeErr

	err := cmd.Start()

	if err != nil {
		DieLoudly(1, "Failed to execute `%s`; error: %s", cmd, err)
	}

	err = cmd.Wait()

	return err == nil
}
