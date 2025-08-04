package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/ecs"
)

func TestRootView_RendersPlanetPanelsThroughGlobalRenderer(t *testing.T) {
	r := NewRootView()
	_, _ = r.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	gs := r.global
	// start new game -> space
	s := NewStartScreen()
	s.result = &StartResult{Action: ActionNewGame}
	gs.subScreen = s
	m, _ := gs.handleStartScreenResult(s.result)
	g := m.(*GlobalScreen)
	// Transitions are completed immediately
	space := g.subScreen.(*SpaceScreen)
	_, _ = space.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	// switch layer to planet and process
	ctx := ecs.GetWorldContext(space.model.World())
	ctx.CurrentLayer = ecs.LayerPlanetSurface
	ecs.SetWorldContext(space.model.World(), ctx)
	_, _ = g.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	// Transition completed automatically
	_, _ = r.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	out := r.View()
	if out == "" {
		t.Fatal("no output")
	}
}
