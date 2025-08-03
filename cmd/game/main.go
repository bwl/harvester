package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/internal/ui"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	mp := ui.NewModelWithRNG(r)
	p := tea.NewProgram(&mp, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
