#!/usr/bin/env -S bash -euo pipefail

projectName="superfile"
version="v1.6.0"
osList=("darwin" "linux" "windows")
archList=("amd64" "arm64")
mkdir dist

# Prevent macOS from adding ._* files to archives
export COPYFILE_DISABLE=1

build_binary() {
    local os="$1"
    local arch="$2"
    local output="$3"

    if [ "$os" = "darwin" ]; then
        env GOOS="$os" GOARCH="$arch" CGO_ENABLED=1 go build -o "$output" main.go
    else
        env GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 go build -o "$output" main.go
    fi
}

for os in "${osList[@]}"; do
    if [ "$os" = "windows" ]; then
        for arch in "${archList[@]}"; do
            echo "$projectName-$os-$version-$arch"
            mkdir "./dist/$projectName-$os-$version-$arch"
            cd ../ || exit
            build_binary "$os" "$arch" "./release/dist/$projectName-$os-$version-$arch/spf.exe"
            cd ./release || exit
            zip -r "./dist/$projectName-$os-$version-$arch.zip" "./dist/$projectName-$os-$version-$arch"
        done
    else
        for arch in "${archList[@]}"; do
            echo "$projectName-$os-$version-$arch"
            mkdir "./dist/$projectName-$os-$version-$arch"
            cd ../ || exit
            build_binary "$os" "$arch" "./release/dist/$projectName-$os-$version-$arch/spf"
            cd ./release || exit
            tar czf "./dist/$projectName-$os-$version-$arch.tar.gz" "./dist/$projectName-$os-$version-$arch"
        done
    fi
done
