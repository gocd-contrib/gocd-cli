package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestExec(t *testing.T) {
	as := asserts(t)
	out := testOut()
	Exec(mockCmd(`echo`, `-n`, `a`, `b`, `c`), testIn(), out, testOut())
	as.eq(`echo -n a b c`, out.String())
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	fmt.Fprint(os.Stdout, strings.Join(os.Args[3:], ` `))
	os.Exit(0)
}

func mockCmd(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func testIn() io.Reader {
	return strings.NewReader(``)
}

func testOut() *strings.Builder {
	return &strings.Builder{}
}
