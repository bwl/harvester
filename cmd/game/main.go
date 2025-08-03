package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/internal/ui"
	"harvester/pkg/ecs"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// load latest autosave if present
	mp := ui.NewModelWithRNG(r)
	if b, err := os.ReadFile(".saves/autosave.gz"); err == nil {
		if s, err := ecs.DecodeSnapshot(b, ecs.SaveOptions{Compress: true}); err == nil {
			_ = ecs.Load(mp.World(), s, nil)
		}
	}
	p := tea.NewProgram(&mp, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
