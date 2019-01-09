package utils

import (
	"io"
	"os"
	"os/exec"
)

func ExecQ(cmd *exec.Cmd) bool {
	return Exec(cmd, os.Stdin, StdoutOrDevNull(), StderrOrDevNull())
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
