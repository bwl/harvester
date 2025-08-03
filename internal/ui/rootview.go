package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/rendering"
)

type RootView struct {
	global *GlobalScreen
	r      *rendering.ViewRenderer
	w, h   int
}

func NewRootView() *RootView { return &RootView{global: NewGlobalScreen()} }

func (r *RootView) Init() tea.Cmd { return r.global.Init() }

func (r *RootView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if wm, ok := msg.(tea.WindowSizeMsg); ok {
		r.w, r.h = wm.Width, wm.Height
		if r.r == nil {
			r.r = rendering.NewViewRenderer(r.w, r.h)
		} else {
			r.r.SetDimensions(r.w, r.h)
		}
	}
	if km, ok := msg.(tea.KeyMsg); ok {
		if km.String() == "ctrl+q" {
			if ih, ok := any(r.global).(InputHandler); ok {
				return r, ih.HandleInput(InputAction{Kind: InputQuit})
			}
		}
		a := mapKeyToAction(km)
		if a.Kind != InputNone {
			if ih, ok := any(r.global).(InputHandler); ok {
				return r, ih.HandleInput(a)
			}
		}
	}
	m, cmd := r.global.Update(msg)
	if g, ok := m.(*GlobalScreen); ok {
		r.global = g
	}
	return r, cmd
}

func mapKeyToAction(k tea.KeyMsg) InputAction {
	s := k.String()
	switch s {
	case "ctrl+c", "ctrl+q":
		return InputAction{Kind: InputQuit}
	case "enter":
		return InputAction{Kind: InputMenuSelect}
	case "esc":
		return InputAction{Kind: InputMenuBack}
	case "w":
		return InputAction{Kind: InputMoveUp}
	case "s":
		return InputAction{Kind: InputMoveDown}
	case "a":
		return InputAction{Kind: InputMoveLeft}
	case "d":
		return InputAction{Kind: InputMoveRight}
	case "h":
		return InputAction{Kind: InputMenuLeft}
	case "l":
		return InputAction{Kind: InputMenuRight}
	case "k":
		return InputAction{Kind: InputMenuUp}
	case "j":
		return InputAction{Kind: InputMenuDown}
	case ">":
		return InputAction{Kind: InputEnter}
	case "ctrl+s":
		return InputAction{Kind: InputSaveAuto}
	case "ctrl+shift+s":
		return InputAction{Kind: InputSaveCompressed}
	case "1":
		return InputAction{Kind: InputSaveSlot1}
	case "2":
		return InputAction{Kind: InputSaveSlot2}
	case "3":
		return InputAction{Kind: InputSaveSlot3}
	}
	return InputAction{Kind: InputNone}
}

func (r *RootView) View() string {
	if r.w == 0 || r.h == 0 {
		r.w, r.h = 80, 24
	}
	if r.r == nil {
		r.r = rendering.NewViewRenderer(r.w, r.h)
	}
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
