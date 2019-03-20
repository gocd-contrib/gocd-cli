#!/bin/bash

set -e

cd "$(dirname "$0")"

PROGNAME="gocd"

rm -f "$PROGNAME"
rm -rf dist

RELEASE="X.x.x"

for arg in $@; do
  case $arg in
    --verbose)
      extra_flags="$extra_flags -v"
      shift
      ;;
    --skip-tests)
      skip=true
      shift
      ;;
    --prod)
      multiplatform=true
      shift
      ;;
    --release)
      RELEASE="$2"
      shift
      ;;
    --release=*)
      RELEASE="${arg#*=}"
      shift
      ;;
    *)
      shift
      ;;
  esac
done

RELEASE="${RELEASE}-${GO_PIPELINE_LABEL:-localbuild}"

function ldflags {
  local _os="${1:-$(go env GOOS)}"
  local _arch="${2:-$(go env GOARCH)}"

  echo "-X main.Version=${RELEASE} -X main.GitCommit=${GIT_COMMIT} -X main.Platform=${_arch}-${_os}"
}

echo "Fetching dependencies"
go get -d $extra_flags ./...

if [[ "true" = "$skip" ]]; then
  echo "Skipping tests"
else
   go test $extra_flags ./...
fi

if (which git &> /dev/null); then
  GIT_COMMIT=$(git rev-list --abbrev-commit -1 HEAD)
else
  GIT_COMMIT="unknown"
fi

if [[ "true" = "$multiplatform" ]]; then
  platforms=(
    darwin/amd64
    freebsd/386
    freebsd/amd64
    linux/386
    linux/amd64
    linux/arm
    linux/arm64
    windows/386
    windows/amd64
  )

  echo "Release: $RELEASE, Revision: $GIT_COMMIT"

  for plt in "${platforms[@]}"; do
    mkdir -p "dist/${plt}"
    arr=(${plt//\// })
    _os="${arr[0]}"
    _arch="${arr[1]}"
    name="$PROGNAME"

    if [ "windows" = "${_os}" ]; then
      name="$name.exe"
    fi

    echo "Building $plt..."

    GOOS="${_os}" go get -d $extra_flags ./...
    GOOS="${_os}" GOARCH="${_arch}" go build \
      -o "dist/${plt}/${name}" \
      -ldflags "$(ldflags "$_os" "$_arch")"
  done
else
  go build \
    -ldflags "$(ldflags)" \
    -o "$PROGNAME"
fi
