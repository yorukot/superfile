#!/usr/bin/env -S bash -euo pipefail

projectName="superfile"
version="v1.1.2"
osList=("darwin" "linux")
archList=("amd64" "arm64")
mkdir dist

for os in "${osList[@]}"; do
    for arch in "${archList[@]}"; do
        echo "$projectName-$os-$version-$arch"
        mkdir "./dist/$projectName-$os-$version-$arch"
        cd ../src || exit
        env GOOS="$os" GOARCH="$arch" go build -o "../release/dist/$projectName-$os-$version-$arch/spf" main.go
        cd ../release || exit
        tar czf "./dist/$projectName-$os-$version-$arch.tar.gz" "./dist/$projectName-$os-$version-$arch"
    done
done
