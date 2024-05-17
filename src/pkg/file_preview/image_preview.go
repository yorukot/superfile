package filepreview

import (
	"encoding/base64"
	"fmt"
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

func imagePreview() string {
    currentTerminal := detectCurrentTerminal()
	if currentTerminal == Kitty {

    } else {
        //...
    }
    return ""
}

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

func showImageInKitty(imagePath string) string {
    imageData, err := pathToBase64Encode(imagePath)
    if err != nil {
        fmt.Println("Error:", err)
        return ""
    }

    encodedData := base64.StdEncoding.EncodeToString([]byte(imageData))
    fmt.Println("\x1b_G=;base64,%s\a", encodedData)
    return fmt.Sprintf("\x1b_G=;base64,%s\a", encodedData)
}