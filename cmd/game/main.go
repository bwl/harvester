package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/internal/ui"
	"harvester/pkg/debug"
)

func main() {
	// Initialize debugging system
	debug.Info("main", "Starting Bubble Rouge game")
	debug.Info("main", "Debug panel available - press F12 to toggle")

	// Create global screen manager that handles all screens and effects
	root := ui.NewRootView()

	// Launch the application with root view
	program := tea.NewProgram(root, tea.WithAltScreen())
	if err := program.Start(); err != nil {
		debug.Errorf("main", "Application failed to start: %v", err)
		fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
		os.Exit(1)
	}

	debug.Info("main", "Game shutting down")
}
