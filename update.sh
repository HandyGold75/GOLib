#!/bin/bash
dir="$( cd "$( dirname "$0" )" && pwd )"

cd "$dir" || exit 1

for subdir in `ls -d */`; do
    if [ "$(ls "$dir/$subdir" | grep go.mod)" = "" ]; then
        continue
    fi


    cd "$dir/$subdir" || exit 1

    file=${subdir,,}
    file=${file/"/"/""}

    go mod edit -go `go version | { read _ _ v _; echo ${v#go}; }`
    go mod tidy
    go get -u ./...
done
