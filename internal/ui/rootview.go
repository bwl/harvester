package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/rendering"
)

type RootView struct{
	global *GlobalScreen
	r *rendering.ViewRenderer
	w,h int
}

func NewRootView() *RootView { return &RootView{ global: NewGlobalScreen() } }

func (r *RootView) Init() tea.Cmd { return r.global.Init() }

func (r *RootView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if wm, ok := msg.(tea.WindowSizeMsg); ok {
		r.w, r.h = wm.Width, wm.Height
		if r.r == nil { r.r = rendering.NewViewRenderer(r.w, r.h) } else { r.r.SetDimensions(r.w, r.h) }
	}
	m, cmd := r.global.Update(msg)
	if g, ok := m.(*GlobalScreen); ok { r.global = g }
	return r, cmd
}

func (r *RootView) View() string {
	if r.w == 0 || r.h == 0 { r.w, r.h = 80, 24 }
	if r.r == nil { r.r = rendering.NewViewRenderer(r.w, r.h) }
	r.r.UnregisterAll()
	
	// Register background layer
	r.r.RegisterContent(newBackgroundLayer(r.w, r.h))
	
	// Let GlobalScreen register its content and effects
	var renderableGlobal RenderableScreen = r.global
	renderableGlobal.RegisterContent(r.r)
	
	// Register TV frame on top
	r.r.RegisterContent(newTVFrame(r.w, r.h))
	
	return r.r.Render()
}
