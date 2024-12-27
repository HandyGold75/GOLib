#!/bin/bash

updatework(){
    go mod edit -go "$(go version | { read -r _ _ v _; echo "${v#go}"; })" || { echo -e "\033[31mFailed: $1.*\033[0m" ; return 1; }
    go mod tidy || { echo -e "\033[31mFailed: $1.*\033[0m" ; return 1; }
    go get -u ./... || { echo -e "\033[31mFailed: $1.*\033[0m" ; return 1; }
}

dir="$( cd "$( dirname "$0" )" && pwd )"
cd "$dir" || exit 1

for subdir in *"/"; do
    if [ ! -f "$dir/$subdir"go.mod ]; then
        continue
    fi

    file=${subdir,,}
    file=${file/"/"/""}
    cd "$dir/$subdir" || { echo -e "\033[31mFailed: $file.bin\033[0m" ; continue; }

    updatework "$file" && echo -e "\033[32mUpdated: $file.bin\033[0m"
done
