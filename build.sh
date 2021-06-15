#!/bin/bash
set -euo pipefail
scriptdir="$(readlink -f "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"
: "$scriptdir"
program_name=quandong
git describe

export CGO_ENABLED=0 # For static build without 'C'

function linux_compile {
  : "$1" # e.g. amd64
  export GOOS=linux
  export GOARCH="$1"
  echo "$GOOS" "$GOARCH"
  rm -f "${program_name}""$GOOS"_"$GOARCH".tgz
  go build -a -ldflags="-X 'main.Version=$(git describe)'" -o "${program_name}""$GOOS"_"$GOARCH" frangipanni.go
  tar zcf "${program_name}""$GOOS"_"$GOARCH".tgz "${program_name}""$GOOS"_"$GOARCH" ./*.lua
}

function cross_compile {
  : "$1" # e.g. linux/amd64
  export GOOS="${1%/*}"
  export GOARCH="${1#*/}"         # This is why bash is so awful
  echo "$GOOS" "$GOARCH"
  rm -f "${program_name}""$GOOS"_"$GOARCH".zip
  go build -a -ldflags="-X 'main.Version=$(git describe)'" -o "${program_name}""$GOOS"_"$GOARCH" frangipanni.go
  zip --quiet "${program_name}""$GOOS"_"$GOARCH".zip "${program_name}""$GOOS"_"$GOARCH" ./*.lua
}

test/confidence.sh

for arch in 386 arm64 amd64
do
  linux_compile "$arch"
done
for dist in windows/amd64 windows/386 darwin/amd64 freebsd/amd64 js/wasm netbsd/amd64 openbsd/amd64
do
  cross_compile "$dist"
done
