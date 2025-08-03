package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

func setupSpaceScreen(t *testing.T) *SpaceScreen {
	gs := NewGlobalScreen()
	s := NewStartScreen()
	s.result = &StartResult{Action: ActionNewGame}
	gs.subScreen = s
	m, _ := gs.handleStartScreenResult(s.result)
	g, ok := m.(*GlobalScreen)
	if !ok {
		t.Fatal("expected GlobalScreen")
	}
	// Transitions are completed immediately
	space, ok := g.subScreen.(*SpaceScreen)
	if !ok {
		t.Fatal("expected SpaceScreen")
	}
	_, _ = space.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	return space
}

func TestSpaceScreen_Movement_WASD(t *testing.T) {
	space := setupSpaceScreen(t)
	m := space.model
	ctx := ecs.GetWorldContext((*m).World())
	ctx.CurrentLayer = ecs.LayerSpace
	ecs.SetWorldContext((*m).World(), ctx)
	pos0, _ := ecs.Get[components.Position]((*m).world, (*m).player)
	space.model.ApplyAction(InputAction{Kind: InputMoveUp})
	for i := 0; i < 60; i++ {
		(*m).scheduler.Update(1.0/20.0, (*m).world)
	}
	pos1, _ := ecs.Get[components.Position]((*m).world, (*m).player)
	if pos1.X == pos0.X && pos1.Y == pos0.Y {
		t.Fatalf("expected movement, got none: (%.2f,%.2f) -> (%.2f,%.2f)", pos0.X, pos0.Y, pos1.X, pos1.Y)
	}
}

func TestSpaceScreen_Movement_Arrows(t *testing.T) {
	space := setupSpaceScreen(t)
	m := space.model
	ctx := ecs.GetWorldContext((*m).World())
	ctx.CurrentLayer = ecs.LayerSpace
	ecs.SetWorldContext((*m).World(), ctx)
	pos0, _ := ecs.Get[components.Position]((*m).world, (*m).player)
	space.model.ApplyAction(InputAction{Kind: InputMoveUp})
	for i := 0; i < 60; i++ {
		(*m).scheduler.Update(1.0/20.0, (*m).world)
	}
	pos1, _ := ecs.Get[components.Position]((*m).world, (*m).player)
	if pos1.X == pos0.X && pos1.Y == pos0.Y {
		t.Fatalf("expected movement with thrust, got none: (%.2f,%.2f) -> (%.2f,%.2f)", pos0.X, pos0.Y, pos1.X, pos1.Y)
	}
}
