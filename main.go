package main

import (
	"embed"
	"github.com/yorukot/superfile/src/cmd"
)

var (
	//go:embed src/superfile_config/*
	content embed.FS
)

func main() {
	cmd.Run(content)
}
