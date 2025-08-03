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
	g.completeTransition()
	// ensure wrapper
	wrapper, ok := g.subScreen.(*GameScreenWrapper)
	if !ok {
		t.Fatal("expected GameScreenWrapper")
	}
	_, _ = wrapper.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	out := wrapper.View()
	if out == "" {
		t.Fatal("no output")
	}
	t.Log(out)
	// HUD assertions
	if !strings.Contains(out, "Fuel") {
		t.Error("HUD missing Fuel")
	}
	if !strings.Contains(out, "Layer") {
		t.Error("HUD missing Layer")
	}
	// Viewport assertions: expect at least height lines
	lines := strings.Split(out, "\n")
	if len(lines) < 24 {
		t.Errorf("expected at least 24 lines, got %d", len(lines))
	}
}
