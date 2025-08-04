package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/rendering"
)

// PlanetScreen handles planet surface and deep exploration
type PlanetScreen struct {
	model         *Model
	width, height int
}

func (p *PlanetScreen) RegisterContent(renderer *rendering.ViewRenderer) {
	if p.width == 0 || p.height == 0 {
		return
	}
	mapH := p.height - 3
	if mapH < 1 {
		mapH = 1
	}
	gm := buildGameGlyphs(p.model, p.width, p.height-3)
	if gm != nil {
		renderer.RegisterContent(newPlanetSurfaceContent(gm))
	}
	renderer.RegisterContent(newHUDContent(p.model))
	renderer.RegisterContent(newUIRightPanel(p.model))
}

func NewPlanetScreen(model *Model) *PlanetScreen {
	return &PlanetScreen{
		model: model,
	}
}

func (p *PlanetScreen) Init() tea.Cmd {
	return p.model.Init()
}

func (p *PlanetScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if wm, ok := msg.(tea.WindowSizeMsg); ok {
		p.width, p.height = wm.Width, wm.Height
	}
	_, cmd := p.model.Update(msg)
	return p, cmd
}

func (p *PlanetScreen) View() string {
	return ""
}

func (p *PlanetScreen) HandleInput(a InputAction) tea.Cmd {
	p.model.ApplyAction(a)
	return nil
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

// SetDimensions implements ResizableScreen interface
func (p *PlanetScreen) SetDimensions(width, height int) {
	p.width = width
	p.height = height
}

// Planet-specific rendering functions


// Use existing itoa function from components.go
