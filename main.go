package main

import (
	"math/rand"
	"time"

	"github.com/gocd-contrib/gocd-cli/cmd"
	"github.com/gocd-contrib/gocd-cli/meta"
)

// These should be set by the linker at build time
var (
	Version   = `devbuild`
	GitCommit = `unknown`
	Platform  = `devbuild`
)

func main() {
	meta.Version = Version
	meta.GitCommit = GitCommit
	meta.Platform = Platform
	cmd.RootCmd.Version = meta.VersionString()

	cmd.Execute()
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
