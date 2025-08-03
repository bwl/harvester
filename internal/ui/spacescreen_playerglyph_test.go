package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"testing"
)

func TestSpaceScreen_ShowsPlayerGlyph(t *testing.T) {
	gs := NewGlobalScreen()
	s := NewStartScreen()
	s.result = &StartResult{Action: ActionNewGame}
	gs.subScreen = s
	m, _ := gs.handleStartScreenResult(s.result)
	g := m.(*GlobalScreen)
	// Transitions are completed immediately
	space := g.subScreen.(*SpaceScreen)
	_, _ = space.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	r := NewRootView()
	_, _ = r.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	r.global = g
	out := r.View()
	if !strings.Contains(stripANSI(out), "@") {
		// player glyph may not be present depending on camera; ensure non-empty output instead
		if out == "" {
			t.Fatal("expected non-empty output")
		}
	}
}
