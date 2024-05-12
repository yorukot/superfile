package filepreview

import (
	"os"
	"strings"
)

const (
    Kitty terminalType = iota
    Konsole
    ITerm2
    WezTerm
    Mintty
    Foot
    Ghostty
    BlackBox
    VSCode
    Tabby
    Hyper
    Unknown
)

type terminalType int

func detectCurrentTerminal() terminalType {
    term := os.Getenv("TERM")
    switch strings.ToLower(term) {
    case "xterm-kitty":
        return Kitty
    case "konsole":
        return Konsole
    case "xterm-256color":
        return ITerm2
    case "wezterm":
        return WezTerm
    case "mintty":
        return Mintty
    case "foot":
        return Foot
    case "ghostty":
        return Ghostty
    case "blackbox":
        return BlackBox
    case "visualstudio":
        return VSCode
    case "tabby":
        return Tabby
    case "hyper":
        return Hyper
    default:
        return Unknown
    }
}