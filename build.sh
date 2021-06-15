#!/bin/bash
set -euo pipefail
scriptdir="$(readlink -f "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"
: "$scriptdir"
program_name=quandong

mkdir -p package bin

git describe

export CGO_ENABLED=0 # For static build without 'C'

function linux_compile {
  : "$1" # e.g. amd64
  export GOOS=linux
  export GOARCH="$1"
  echo "$GOOS" "$GOARCH"
  rm -f package/"${program_name}""$GOOS"_"$GOARCH".tgz
  go build -a -ldflags="-X 'main.Version=$(git describe)'" -o bin/"${program_name}"_"$GOOS"_"$GOARCH" "${program_name}".go
  tar zcf package/"${program_name}"_"$GOOS"_"$GOARCH".tgz bin/"${program_name}"_"$GOOS"_"$GOARCH"
}

function cross_compile {
  : "$1" # e.g. linux/amd64
  export GOOS="${1%/*}"
  export GOARCH="${1#*/}"         # This is why bash is so awful
  echo "$GOOS" "$GOARCH"
  rm -f "${program_name}""$GOOS"_"$GOARCH".zip
  go build -a -ldflags="-X 'main.Version=$(git describe)'" -o bin/"${program_name}"_"$GOOS"_"$GOARCH" "${program_name}".go
  zip --quiet package/"${program_name}"_"$GOOS"_"$GOARCH".zip bin/"${program_name}"_"$GOOS"_"$GOARCH"
}

linux_compile amd64
cd bin
ln -sf quandong_linux_amd64 "$PWD"/quandong
ln -sf quandong_linux_amd64 "$PWD"/ls
exit 0
for arch in 386 arm64 amd64
do
  linux_compile "$arch"
done
for dist in windows/amd64 windows/386 darwin/amd64 freebsd/amd64 js/wasm netbsd/amd64 openbsd/amd64
do
  cross_compile "$dist"
done
