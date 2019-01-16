#!/bin/bash

set -e

rm -f gocd gocd.exe
rm -rf build

RELEASE="UNSPECIFIED"

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

echo "Fetching dependencies"
go get -d $extra_flags ./...

echo "Fetching any windows-specific dependencies"
GOOS="windows" go get -d $extra_flags ./... # get any windows-specific deps as well

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
    mkdir -p "build/$plt"
    arr=(${plt//\// })
    _os="${arr[0]}"
    _arch="${arr[1]}"
    name="gocd"

    if [[ "windows" = "${_os}" ]]; then
      name="$name.exe"
    fi

    echo "Building $plt..."

    GOOS="${_os}" GOARCH="${_arch}" go build \
      -o "build/${plt}/${name}" \
      -ldflags "-X main.Version=$RELEASE -X main.GitCommit=$GIT_COMMIT -X main.Platform=$_arch-$_os" \
      main.go
  done
else
  _arch=$(go env GOARCH)
  _os=$(go env GOOS)

  go build \
    -ldflags "-X main.Version=X.x.x-devbuild -X main.GitCommit=$GIT_COMMIT -X main.Platform=$_arch-$_os" \
    -o gocd \
    main.go
fi
