package main

import (
	"fmt"
	"os"

	"uldocker/internal/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(app.NewModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
