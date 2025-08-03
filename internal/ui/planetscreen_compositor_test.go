package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/ecs"
	"harvester/pkg/rendering"
)

func TestPlanetScreen_CompositorRendersHUDAndMap(t *testing.T) {
	gs := NewGlobalScreen()
	// transition to game (space)
	s := NewStartScreen()
	s.result = &StartResult{Action: ActionNewGame}
	gs.subScreen = s
	m, _ := gs.handleStartScreenResult(s.result)
	g := m.(*GlobalScreen)
	// Transitions are completed immediately, move to planet layer to trigger transition
	space := g.subScreen.(*SpaceScreen)
	_, _ = space.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	ctx := ecs.GetWorldContext(space.model.World())
	ctx.CurrentLayer = ecs.LayerPlanetSurface
	ecs.SetWorldContext(space.model.World(), ctx)
	_, _ = g.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	// Transitions are completed immediately
	planet, ok := g.subScreen.(*PlanetScreen)
	if !ok {
		t.Fatal("expected PlanetScreen")
	}
	_, _ = planet.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	vr := rendering.NewViewRenderer(80, 24)
	var rs RenderableScreen = planet
	rs.RegisterContent(vr)
	out := vr.Render()
	if out == "" {
		t.Fatal("no output")
	}
	if !strings.Contains(out, "HP:") || !strings.Contains(out, "Layer") {
		t.Error("HUD missing expected fields")
	}
}
