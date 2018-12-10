package utils

import (
	"log"
	"os"
	"os/exec"
)

func Exec(cmd *exec.Cmd, pipeIn *os.File, pipeOut *os.File, pipeErr *os.File) bool {
	cmd.Stdin = pipeIn
	cmd.Stdout = pipeOut
	cmd.Stderr = pipeErr

	err := cmd.Start()

	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	err = cmd.Wait()

	return err == nil
}
