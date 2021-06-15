#!/bin/bash
set -euo pipefail
scriptdir="$(readlink -f "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"
: "$scriptdir"
program_name=quandong

mkdir -p bin

export CGO_ENABLED=0 # For static build without 'C'
export GOOS=linux
export GOARCH=amd64
go build -a -ldflags="-X 'main.Version=$(git describe)'" -o bin/"${program_name}"_"$GOOS"_"$GOARCH" "${program_name}".go

mv bin/quandong_linux_amd64 bin/quandong
(cd bin; ln -sf quandong "$PWD"/date; cp quandong uname)
export PATH=$PWD/bin:$PATH
hash -r
rm -f quandong-date*json quandong-uname*json
which date
which uname
date +%s
uname -a
ls -ltr quandong-{date,uname}*json