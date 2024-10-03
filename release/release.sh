#!/usr/bin/env -S bash -euo pipefail

projectName="superfile"
version="v1.1.5"
osList=("darwin" "linux" "windows")
archList=("amd64" "arm64")
mkdir dist

for os in "${osList[@]}"; do
    if [ "$os" = "windows" ]; then
        for arch in "${archList[@]}"; do
            echo "$projectName-$os-$version-$arch"
            mkdir "./dist/$projectName-$os-$version-$arch"
            cd ../ || exit
            env GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 go build -o "./release/dist/$projectName-$os-$version-$arch/spf.exe" main.go
            cd ./release || exit
            zip -r "./dist/$projectName-$os-$version-$arch.zip" "./dist/$projectName-$os-$version-$arch"
        done
    else
        for arch in "${archList[@]}"; do
            echo "$projectName-$os-$version-$arch"
            mkdir "./dist/$projectName-$os-$version-$arch"
            cd ../ || exit
            env GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 go build -o "./release/dist/$projectName-$os-$version-$arch/spf" main.go
            cd ./release || exit
            tar czf "./dist/$projectName-$os-$version-$arch.tar.gz" "./dist/$projectName-$os-$version-$arch"
        done
    fi
done
w