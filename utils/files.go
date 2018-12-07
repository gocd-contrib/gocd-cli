package utils

import (
	"os"
)

func IsFile(name string) bool {
	fi, err := os.Stat(name)
	return err == nil && fi.Mode().IsRegular()
}

func IsDir(name string) bool {
	fi, err := os.Stat(name)
	return err == nil && fi.IsDir()
}
