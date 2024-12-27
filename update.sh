#!/bin/bash

update(){
    go get go@latest || { echo -e "\033[31mFailed: $1.*\033[0m" ; return 1; }
    go get -u || { echo -e "\033[31mFailed: $1.*\033[0m" ; return 1; }
    go mod tidy || { echo -e "\033[31mFailed: $1.*\033[0m" ; return 1; }
    for indirect in $(tail +3 go.mod | grep "// indirect" | awk '{if ($1 =="require") print $2; else print $1;}'); do go get -u "${indirect}"; done
    go get -u || { echo -e "\033[31mFailed: $1.*\033[0m" ; return 1; }
    go mod tidy || { echo -e "\033[31mFailed: $1.*\033[0m" ; return 1; }
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

    update "$file" && echo -e "\033[32mUpdated: $file.bin\033[0m"
done
