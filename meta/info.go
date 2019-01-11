package meta

import "fmt"

// populated by main.go upon execute
var Version string
var GitCommit string
var Platform string

func VersionString() string {
	return fmt.Sprintf(`%s (%s rev. %s)`, Version, Platform, GitCommit)
}
