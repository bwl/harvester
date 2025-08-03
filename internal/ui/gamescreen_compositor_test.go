package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestGameScreen_CompositorRendersHUDAndMap(t *testing.T) {
	gs := NewGlobalScreen()
	// transition to game
	s := NewStartScreen()
	s.result = &StartResult{Action: ActionNewGame}
	gs.subScreen = s
	m, _ := gs.handleStartScreenResult(s.result)
	g, ok := m.(*GlobalScreen)
	if !ok {
		t.Fatal("expected GlobalScreen")
	}
	// ensure space screen (games start in space)
	spaceScreen, ok := g.subScreen.(*SpaceScreen)
	if !ok {
		t.Fatal("expected SpaceScreen")
	}
	_, _ = spaceScreen.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	// Use RootView to render via global compositor
	r := NewRootView()
	_, _ = r.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	r.global = g
	out := r.View()
	if out == "" {
		t.Fatal("no output")
	}
	lines := strings.Split(out, "\n")
	if len(lines) < 24 {
		t.Errorf("expected at least 24 lines, got %d", len(lines))
	}
}
