package ui

import "github.com/charmbracelet/lipgloss"

// Layout represents the overall UI layout configuration
type Layout struct {
	Width  int
	Height int
	Margin int

	// Panel dimensions
	RightPanelWidth   int
	BottomPanelHeight int
	MinMapWidth       int
	MinMapHeight      int
}

// NewLayout creates a layout with default values
func NewLayout(width, height int) Layout {
	return Layout{
		Width:             width,
		Height:            height,
		Margin:            2,
		RightPanelWidth:   30,
		BottomPanelHeight: 5,
		MinMapWidth:       10,
		MinMapHeight:      10,
	}
}

// LayoutDimensions calculates the dimensions for each panel
type LayoutDimensions struct {
	ContentWidth  int
	ContentHeight int
	MapWidth      int
	MapHeight     int
	RightWidth    int
	RightHeight   int
	LogWidth      int
	LogHeight     int
}

// Calculate returns the calculated dimensions for all panels
func (l Layout) Calculate() LayoutDimensions {
	// Account for margin
	contentW := l.Width - (l.Margin * 2)
	contentH := l.Height - (l.Margin * 2)

	// Calculate map dimensions
	mapW := contentW - l.RightPanelWidth
	mapH := contentH - l.BottomPanelHeight

	// Ensure minimum sizes
	if mapW < l.MinMapWidth {
		mapW = l.MinMapWidth
	}
	if mapH < l.MinMapHeight {
		mapH = l.MinMapHeight
	}

	return LayoutDimensions{
		ContentWidth:  contentW,
		ContentHeight: contentH,
		MapWidth:      mapW,
		MapHeight:     mapH,
		RightWidth:    l.RightPanelWidth,
		RightHeight:   mapH, // Right panel same height as map
		LogWidth:      contentW,
		LogHeight:     l.BottomPanelHeight,
	}
}

// LayoutPreset represents different layout configurations
type LayoutPreset int

const (
	LayoutFull LayoutPreset = iota
	LayoutCompact
	LayoutMobile
)

// ApplyPreset modifies the layout based on the preset
func (l *Layout) ApplyPreset(preset LayoutPreset) {
	switch preset {
	case LayoutFull:
		l.RightPanelWidth = 30
		l.BottomPanelHeight = 5
		l.Margin = 2
	case LayoutCompact:
		l.RightPanelWidth = 25
		l.BottomPanelHeight = 3
		l.Margin = 1
	case LayoutMobile:
		l.RightPanelWidth = 20
		l.BottomPanelHeight = 2
		l.Margin = 0
	}
}

// Validate ensures the layout dimensions are reasonable
func (l Layout) Validate() bool {
	dims := l.Calculate()
	return dims.MapWidth >= l.MinMapWidth &&
		dims.MapHeight >= l.MinMapHeight &&
		dims.ContentWidth > 0 &&
		dims.ContentHeight > 0
}

// LayoutManager handles responsive layout adjustments
type LayoutManager struct {
	currentLayout Layout
	autoResize    bool
}

// NewLayoutManager creates a new layout manager
func NewLayoutManager(width, height int) *LayoutManager {
	return &LayoutManager{
		currentLayout: NewLayout(width, height),
		autoResize:    true,
	}
}

// Update updates the layout dimensions and applies responsive adjustments
func (lm *LayoutManager) Update(width, height int) {
	lm.currentLayout.Width = width
	lm.currentLayout.Height = height

	if lm.autoResize {
		lm.applyResponsiveLayout()
	}
}

// applyResponsiveLayout automatically adjusts layout based on screen size
func (lm *LayoutManager) applyResponsiveLayout() {
	w, h := lm.currentLayout.Width, lm.currentLayout.Height

	// Very small screens
	if w < 80 || h < 20 {
		lm.currentLayout.ApplyPreset(LayoutMobile)
	} else if w < 120 || h < 30 {
		// Small screens
		lm.currentLayout.ApplyPreset(LayoutCompact)
	} else {
		// Normal screens
		lm.currentLayout.ApplyPreset(LayoutFull)
	}
}

// GetLayout returns the current layout
func (lm *LayoutManager) GetLayout() Layout {
	return lm.currentLayout
}

// SetAutoResize enables or disables automatic responsive layout
func (lm *LayoutManager) SetAutoResize(enabled bool) {
	lm.autoResize = enabled
}

// RenderWithLayout renders content using the layout system
func (lm *LayoutManager) RenderWithLayout(mapStr, rightStr, statusStr, logStr string) string {
	dims := lm.currentLayout.Calculate()

	mapPanel := Sized(dims.MapWidth, dims.MapHeight, mapStr)
	rightPanel := Sized(dims.RightWidth, 0, rightStr)

	main := lipgloss.JoinHorizontal(lipgloss.Top, mapPanel, rightPanel)

	content := lipgloss.JoinVertical(lipgloss.Left,
		statusStr,
		main,
		Sized(dims.LogWidth, 0, logStr),
	)

	return lipgloss.NewStyle().
		Margin(lm.currentLayout.Margin).
		Render(content)
}
