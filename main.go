package main

import (
	"embed"

	cmd "github.com/MHNightCat/superfile/src/cmd"
)

var (
	//go:embed src/superfileConfig/*
	content embed.FS
)

func main() {
	cmd.Run(content)
}