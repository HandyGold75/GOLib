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

    # go get -u
    go mod tidy
    go get -u -d ./...
done
