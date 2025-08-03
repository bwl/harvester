package ui

import (
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/ecs"
	"harvester/pkg/rendering"
)

// GameScreenWrapper wraps the existing Model to implement SubScreen interface
type GameScreenWrapper struct {
	model         *Model
	renderer      *rendering.ViewRenderer
	width, height int
}

func NewGameScreenWrapper(startResult *StartResult) *GameScreenWrapper {
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

	return &GameScreenWrapper{
		model: &model,
	}
}

func (g *GameScreenWrapper) Init() tea.Cmd {
	return g.model.Init()
}

func (g *GameScreenWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle game-specific quit with ESC
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "escape" {
			// TODO: Return to start screen or show pause menu
			// For now, just quit with animation
			return g, tea.Quit
		}
	}

	newModel, cmd := g.model.Update(msg)
	if newModel != g.model {
		// Model was replaced, update our wrapper
		if gameModel, ok := newModel.(*Model); ok {
			g.model = gameModel
		}
	}

	return g, cmd
}

func (g *GameScreenWrapper) View() string {
	if g.renderer == nil || g.width == 0 || g.height == 0 {
		return g.model.View()
	}
	g.renderer.UnregisterAll()
	mapH := g.height - 3
	if mapH < 1 {
		mapH = 1
	}
	gm := buildGameGlyphs(g.model, g.width, mapH)
	if gm != nil {
		g.renderer.RegisterContent(newTerrainContent(gm))
	}
	hud := buildHUDGlyphs(g.model, g.width)
	if hud != nil {
		g.renderer.RegisterContent(newHUDContent(hud))
	}
	rp := buildRightPanelGlyphs(g.model)
	if rp != nil {
		g.renderer.RegisterContent(newUIRightPanel(rp))
	}
	return g.renderer.Render()
}

func (g *GameScreenWrapper) HandleGlobalAction(action GlobalAction) (SubScreen, tea.Cmd) {
	switch action {
	case ActionStartShutdown:
		// Game can handle shutdown by saving state
		return g, nil
	default:
		return g, nil
	}
}

// Use existing itoa function from components.go
