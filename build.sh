#!/bin/bash

set -e

go get github.com/mitchellh/go-homedir \
  github.com/spf13/cobra \
  github.com/spf13/viper \
  github.com/blang/semver \
  github.com/dustin/go-humanize

rm -f gocd

for arg in $@; do
  case $arg in
    --skip-tests)
      skip=true
      shift
      ;;
    *)
      shift
      ;;
  esac
done

if [[ "true" = "$skip" ]]; then
  echo "Skipping tests"
else
   go test ./...
fi

go build -o gocd main.go
