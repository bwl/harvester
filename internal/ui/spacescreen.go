package ui

import (
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/ecs"
	"harvester/pkg/rendering"
)

// SpaceScreen handles space navigation and planet selection
type SpaceScreen struct {
	model         *Model
	renderer      *rendering.ViewRenderer
	width, height int
}

func NewSpaceScreen(startResult *StartResult) *SpaceScreen {
	// Create the game model with random seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	model := NewModelWithRNG(r)

	// Load save data based on start result
	switch startResult.Action {
	case ActionContinue:
		// Load autosave
		if b, err := os.ReadFile(".saves/autosave.gz"); err == nil {
			if s, err := ecs.DecodeSnapshot(b, ecs.SaveOptions{Compress: true}); err == nil {
				_ = ecs.Load(model.World(), s, nil)
			}
		}

	case ActionLoadSlot:
		// Load specific slot
		slotPath := ".saves/slot" + itoa(startResult.SlotNum) + ".gz"
		if b, err := os.ReadFile(slotPath); err == nil {
			if s, err := ecs.DecodeSnapshot(b, ecs.SaveOptions{Compress: true}); err == nil {
				_ = ecs.Load(model.World(), s, nil)
			}
		}

	case ActionNewGame:
		// Start fresh - no loading needed, starts in LayerSpace by default
	}

	return &SpaceScreen{
		model: &model,
	}
}

func (s *SpaceScreen) Init() tea.Cmd {
	return s.model.Init()
}

func (s *SpaceScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle space-specific controls
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "escape":
			// Return to start screen
			return s, tea.Quit // TODO: Should transition back to start screen
		case ">":
			// Check if we should transition to planet screen
			if s.shouldTransitionToPlanet() {
				// TODO: Return a message to transition to planet screen
				return s, nil
			}
		}
	}

	// Forward the message to the underlying model
	_, cmd := s.model.Update(msg)
	
	return s, cmd
}

func (s *SpaceScreen) View() string {
	if s.renderer == nil || s.width == 0 || s.height == 0 {
		return s.model.View()
	}
	s.renderer.UnregisterAll()
	
	// Render space-specific content
	mapH := s.height - 3
	if mapH < 1 {
		mapH = 1
	}
	
	// Build space glyphs (stars, planets, player ship)
	gm := buildSpaceGlyphs(s.model, s.width, mapH)
	if gm != nil {
		s.renderer.RegisterContent(newSpaceContent(gm))
	}
	
	// Build space HUD (fuel, coordinates, planet info)
	hud := buildSpaceHUDGlyphs(s.model, s.width)
	if hud != nil {
		s.renderer.RegisterContent(newHUDContent(hud))
	}
	
	return s.renderer.Render()
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

// Check if player is over a planet and pressed enter
func (s *SpaceScreen) shouldTransitionToPlanet() bool {
	ctx := ecs.GetWorldContext(s.model.World())
	return ctx.CurrentLayer == ecs.LayerPlanetSurface
}

// Space-specific rendering functions
func buildSpaceGlyphs(m *Model, width, height int) [][]rendering.Glyph {
	// Filter for space layer content only
	ctx := ecs.GetWorldContext(m.World())
	if ctx.CurrentLayer != ecs.LayerSpace {
		return nil
	}
	
	// Use existing map rendering but focused on space content
	return buildGameGlyphs(m, width, height)
}

func buildSpaceHUDGlyphs(m *Model, width int) [][]rendering.Glyph {
	// Space-specific HUD: fuel gauge, coordinates, planet distances
	return buildHUDGlyphs(m, width)
}

func newSpaceContent(gm [][]rendering.Glyph) rendering.RenderableContent {
	return newTerrainContent(gm)
}