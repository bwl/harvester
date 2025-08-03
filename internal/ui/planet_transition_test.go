package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/ecs"
)

func TestGlobal_Transitions_To_Planet_WhenOnSurface(t *testing.T) {
	gs := NewGlobalScreen()
	// go to space
	s := NewStartScreen()
	s.result = &StartResult{Action: ActionNewGame}
	gs.subScreen = s
	m, _ := gs.handleStartScreenResult(s.result)
	g, ok := m.(*GlobalScreen)
	if !ok {
		t.Fatal("expected GlobalScreen")
	}
	// Transitions are completed immediately, ensure space
	space, ok := g.subScreen.(*SpaceScreen)
	if !ok {
		t.Fatal("expected SpaceScreen")
	}
	_, _ = space.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	// set layer to planet surface
	ctx := ecs.GetWorldContext(space.model.World())
	ctx.CurrentLayer = ecs.LayerPlanetSurface
	ecs.SetWorldContext(space.model.World(), ctx)
	// send any msg to drive GlobalScreen.Update
	_, _ = g.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	// Transition should be completed automatically
	if _, ok := g.subScreen.(*PlanetScreen); !ok {
		t.Fatalf("expected PlanetScreen after transition, got %T", g.subScreen)
	}
}
