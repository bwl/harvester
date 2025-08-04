package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/debug"
	"harvester/pkg/ecs"
	"harvester/pkg/rendering"
)

// SpaceScreen handles space navigation and planet selection
type SpaceScreen struct {
	model  *Model
	width  int
	height int
}

func (s *SpaceScreen) RegisterContent(renderer *rendering.ViewRenderer) {
	w, h := renderer.GetDimensions()

	// Use SpaceScreen's own dimensions if available, fallback to renderer dimensions
	if s.width > 0 && s.height > 0 {
		w, h = s.width, s.height
	}

	if w == 0 || h == 0 {
		debug.Warn("spacescreen", "Zero dimensions received in RegisterContent, skipping render")
		return
	}
	// Build space map via unified render system using current world
	gm := buildGameGlyphs(s.model, w, h-3)
	if gm != nil {
		renderer.RegisterContent(newExpanseContent(gm))
	} else {
		debug.Warn("spacescreen", "buildGameGlyphs returned nil")
	}
}

func NewSpaceScreen(model *Model) *SpaceScreen {
	return &SpaceScreen{
		model: model,
	}
}

func (s *SpaceScreen) Init() tea.Cmd {
	return s.model.Init()
}

func (s *SpaceScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Always forward to model so input and ticks are processed
	if act, ok := msg.(GlobalAction); ok && act == ActionSwitchToGame {
		if s.shouldTransitionToPlanet() {
			return s, tea.Quit
		}
	}
	_, cmd := s.model.Update(msg)
	return s, cmd
}

func (s *SpaceScreen) View() string { return s.model.View() }

func (s *SpaceScreen) HandleInput(a InputAction) tea.Cmd {
	s.model.ApplyAction(a)
	return nil
}

func (s *SpaceScreen) HandleGlobalAction(action GlobalAction) (SubScreen, tea.Cmd) {
	switch action {
	case ActionStartShutdown:
		// Space screen can handle shutdown by saving state
		return s, nil
	default:
		return s, nil
	}
}

// SetDimensions implements ResizableScreen interface
func (s *SpaceScreen) SetDimensions(width, height int) {
	s.width = width
	s.height = height
}

// Check if player is over a planet and pressed enter
func (s *SpaceScreen) shouldTransitionToPlanet() bool {
	ctx := ecs.GetWorldContext(s.model.World())
	return ctx.CurrentLayer == ecs.LayerPlanetSurface
}
