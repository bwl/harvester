package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/debug"
	"harvester/pkg/rendering"
)

type RootView struct {
	global     *GlobalScreen
	r          *rendering.ViewRenderer
	w, h       int
	debugPanel *DebugPanel
	showDebug  bool
}

func NewRootView() *RootView {
	return &RootView{
		global:     NewGlobalScreen(),
		debugPanel: NewDebugPanel(),
	}
}

func (r *RootView) Init() tea.Cmd { return r.global.Init() }

func (r *RootView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if wm, ok := msg.(tea.WindowSizeMsg); ok {
		r.w, r.h = wm.Width, wm.Height
		if r.r == nil {
			r.r = rendering.NewViewRenderer(r.w, r.h)
		} else {
			r.r.SetDimensions(r.w, r.h)
		}
		if r.debugPanel != nil {
			r.debugPanel.SetDimensions(r.w, r.h)
		}
	}
	if km, ok := msg.(tea.KeyMsg); ok {
		// Check for debug panel toggle first
		if km.String() == "f12" {
			r.showDebug = !r.showDebug
			return r, nil
		}

		// If debug panel is open, handle its input
		if r.showDebug && r.debugPanel != nil {
			switch km.String() {
			case "tab":
				r.debugPanel.NextTab()
				return r, nil
			case "shift+tab":
				r.debugPanel.PrevTab()
				return r, nil
			case "up":
				r.debugPanel.ScrollUp()
				return r, nil
			case "down":
				r.debugPanel.ScrollDown()
				return r, nil
			case "f":
				if r.debugPanel.ActiveTab == TabLogs {
					r.debugPanel.CycleLogFilter()
				}
				return r, nil
			case "c":
				if r.debugPanel.ActiveTab == TabLogs {
					debug.Clear()
					return r, nil
				}
			case "esc", "q":
				r.showDebug = false
				return r, nil
			}
		}

		if km.String() == "ctrl+q" {
			if ih, ok := any(r.global).(InputHandler); ok {
				return r, ih.HandleInput(InputAction{Kind: InputQuit})
			}
		}

		// Only process normal game input if debug panel is not showing
		if !r.showDebug {
			a := mapKeyToAction(km)
			if a.Kind != InputNone {
				if ih, ok := any(r.global).(InputHandler); ok {
					return r, ih.HandleInput(a)
				}
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
	case "f12":
		return InputAction{Kind: InputDebugToggle}
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

	baseView := r.r.Render()

	// If debug panel is active, render it on top
	if r.showDebug && r.debugPanel != nil {
		if r.debugPanel.width == 0 || r.debugPanel.height == 0 {
			r.debugPanel.SetDimensions(r.w, r.h)
		}
		return r.debugPanel.Render()
	}

	return baseView
}
