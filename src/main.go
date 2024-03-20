package main

import (
	"fmt"
	"os"
	"superfile/src/components"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(components.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
