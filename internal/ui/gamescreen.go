package ui

import (
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/ecs"
	"harvester/pkg/rendering"
)

// PlanetScreen handles planet surface and deep exploration
type PlanetScreen struct {
	model         *Model
	renderer      *rendering.ViewRenderer
	width, height int
}

func NewPlanetScreen(startResult *StartResult) *PlanetScreen {
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
		// Start fresh - no loading needed
	}

	return &PlanetScreen{
		model: &model,
	}
}

func (p *PlanetScreen) Init() tea.Cmd {
	return p.model.Init()
}

func (p *PlanetScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle planet-specific controls
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "escape" {
			// TODO: Return to space screen when escaping from planet
			return p, tea.Quit
		}
	}

	// Forward the message to the underlying model
	_, cmd := p.model.Update(msg)
	
	return p, cmd
}

func (p *PlanetScreen) View() string {
	if p.renderer == nil || p.width == 0 || p.height == 0 {
		return p.model.View()
	}
	p.renderer.UnregisterAll()
	mapH := p.height - 3
	if mapH < 1 {
		mapH = 1
	}
	
	// Build planet-specific glyphs (terrain, creatures, items)
	gm := buildPlanetGlyphs(p.model, p.width, mapH)
	if gm != nil {
		p.renderer.RegisterContent(newTerrainContent(gm))
	}
	
	// Build planet HUD (health, inventory, depth)
	hud := buildPlanetHUDGlyphs(p.model, p.width)
	if hud != nil {
		p.renderer.RegisterContent(newHUDContent(hud))
	}
	
	rp := buildRightPanelGlyphs(p.model)
	if rp != nil {
		p.renderer.RegisterContent(newUIRightPanel(rp))
	}
	return p.renderer.Render()
}

func (p *PlanetScreen) HandleGlobalAction(action GlobalAction) (SubScreen, tea.Cmd) {
	switch action {
	case ActionStartShutdown:
		// Planet screen can handle shutdown by saving state
		return p, nil
	default:
		return p, nil
	}
}

// Planet-specific rendering functions
func buildPlanetGlyphs(m *Model, width, height int) [][]rendering.Glyph {
	// Filter for planet layer content only
	ctx := ecs.GetWorldContext(m.World())
	if ctx.CurrentLayer == ecs.LayerSpace {
		return nil // Don't render space content on planet screen
	}
	
	// Use existing map rendering but focused on planet content
	return buildGameGlyphs(m, width, height)
}

func buildPlanetHUDGlyphs(m *Model, width int) [][]rendering.Glyph {
	// Planet-specific HUD: health, inventory, depth, temperature, etc.
	return buildHUDGlyphs(m, width)
}

// Use existing itoa function from components.go
