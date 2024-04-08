#!/bin/bash

projectName="superfile"
version="v1.0.1"
osList=("darwin" "freebsd" "linux" "openbsd" "netbsd")
archList=("amd64" "arm" "arm64")
mkdir dist

for os in "${osList[@]}"; do
    for arch in "${archList[@]}"; do
        echo $projectName-$os-$version-$arch
        mkdir ./dist/$projectName-$os-$version-$arch
        cd ../src
        go build -o ../release/dist/$projectName-$os-$version-$arch/spf main.go 
        cd ../release
        tar czf ./dist/$projectName-$os-$version-$arch.tar.gz ./dist/$projectName-$os-$version-$arch
    done
done