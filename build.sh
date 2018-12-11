#!/bin/bash

set -e

go get github.com/mitchellh/go-homedir \
  github.com/spf13/cobra \
  github.com/spf13/viper

rm -f gocd
go build -o gocd main.go
