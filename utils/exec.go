package utils

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
)

func Exec(cmd *exec.Cmd, pipeOut *os.File, pipeErr *os.File) bool {
	var stdoutBuf, stderrBuf bytes.Buffer

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(pipeOut, &stdoutBuf)
	stderr := io.MultiWriter(pipeErr, &stderrBuf)

	err := cmd.Start()

	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}

	return err == nil
}
