package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"bubbleRouge/internal/ui"
	"bubbleRouge/pkg/engine"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	gs := engine.New(time.Now().UnixNano())
	p := tea.NewProgram(ui.NewModel(gs), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
