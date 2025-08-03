package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/internal/ui"
)

func main() {
	// Create global screen manager that handles all screens and effects
	root := ui.NewRootView()

	// Launch the application with root view
	program := tea.NewProgram(root, tea.WithAltScreen())
	if err := program.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
		os.Exit(1)
	}
}
