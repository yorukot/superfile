#!/usr/bin/env bash

# tidy up the go mod
go mod tidy

# format the code
go fmt ./...

# run the linter
golangci-lint run

# build the app
CGO_ENABLED=0 go build -o ./bin/spf
