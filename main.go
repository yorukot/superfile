package main

import (
	"embed"

	cmd "github.com/yorukot/superfile/src/cmd"
)

var (
	//go:embed src/superfileConfig/*
	content embed.FS
)

func main() {
	cmd.Run(content)
}
