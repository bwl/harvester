package ui

import (
	"github.com/charmbracelet/lipgloss/v2"
	"harvester/pkg/rendering"
)

type uiRightPanel struct {
	model *Model
	w, h  int
}

func newUIRightPanel(model *Model) *uiRightPanel {
	return &uiRightPanel{model: model, w: 20, h: 10} // Default size
}

func (u *uiRightPanel) GetLayer() rendering.Layer { return rendering.LayerUI }
func (u *uiRightPanel) GetZ() int                 { return rendering.ZUI }

func (u *uiRightPanel) ToLipglossLayer() *lipgloss.Layer {
	if u.model == nil {
		return lipgloss.NewLayer("").X(0).Y(2).Z(u.GetZ()).ID("right-panel")
	}

	// For now, just return empty panel since buildRightPanelGlyphs returns nil
	content := ""

	style := lipgloss.NewStyle().
		Width(u.w).
		Height(u.h).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#444444"))

	styledContent := style.Render(content)

	return lipgloss.NewLayer(styledContent).
		X(-u.w). // Position from right edge
		Y(2).
		Z(u.GetZ()).
		ID("right-panel")
}
