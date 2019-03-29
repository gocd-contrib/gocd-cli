#!/bin/bash

set -e

function main {
  local progname=$(basename "$0")

  if [ $# -eq 0 ]; then
    die "${progname} requires at least 1 argument"
  fi

  check_for_files signing-keys/codesign.keychain.password.gpg
  check_deps codesign security gpg

  local binary="$1"

  if [ ! -s "$binary" ]; then
    die "Executable [${binary}] could not be found or is empty"
  fi

  decrypt_keychain_passwd codesign.keychain.password

  unlock_codesign_keychain codesign.keychain.password

  codesign_binary "$binary"

  lock_codesign_keychain

  package_binary
}

function package_binary {
  (cd dist/darwin/amd64 && \
    zip -r ../../../osx-cli.zip gocd)
}

function codesign_binary {
  local binary="$1"
  chmod a+rx "$binary"

  echo "SHA-256 before codesign: $(shasum -a 256 -b "$binary")"
  codesign --force --verify --verbose --sign "Developer ID Application: ThoughtWorks (LL62P32G5C)" "$binary"
  echo "SHA-256 after codesign: $(shasum -a 256 -b "$binary")"
}

function unlock_codesign_keychain {
  local passwd_file="$1"

  security unlock-keychain -p "$(cat "$passwd_file")" "$(keychain_path)"
}

function lock_codesign_keychain {
  security lock-keychain "$(keychain_path)"
}

function decrypt_keychain_passwd {
  local outfile="$1"

  if ! gpg --batch --yes --passphrase-file "$(gpg_passwd)" --output "$outfile" signing-keys/codesign.keychain.password.gpg; then
    die "Failed to decrypt codesigning keychain password"
  fi
}

function keychain_path {
  if [ -r "${HOME}/Library/Keychains/codesign.keychain" ]; then
    echo "${HOME}/Library/Keychains/codesign.keychain"
  elif [ -r "${HOME}/Library/Keychains/codesign.keychain-db" ]; then
    echo "${HOME}/Library/Keychains/codesign.keychain-db"
  else
    die "You don't appear to have the codesigning keychain"
  fi
}

function gpg_passwd {
  if [ -n "$GOCD_GPG_PASSPHRASE" ]; then
    echo "$GOCD_GPG_PASSPHRASE" > gpg-passphrase
  fi

  if [ ! -s "gpg-passphrase" ]; then
    die "Cannot codesign without the GPG passphrase! Please set the GOCD_GPG_PASSPHRASE env variable"
  else
    echo "gpg-passphrase"
  fi
}

function check_for_files {
  if [ $# -eq 0 ]; then
    die "check_for_files() requires at least 1 argument"
  fi

  for f in $@; do
    if [ ! -r "$f" ]; then
      die "This script requires the file ${f} to exist and be readable"
    fi
  done
}

function check_deps {
  if [ $# -eq 0 ]; then
    die "check_deps() requires at least 1 argument"
  fi

  for dep in $@; do
    if ! which "$dep" &> /dev/null; then
      die "This script requires ${dep} to be in the PATH"
    fi
  done
}

function die {
  >&2 echo $@
  exit 1
}

main $@
