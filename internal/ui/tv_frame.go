package ui

import (
	"github.com/charmbracelet/lipgloss/v2"
	"harvester/pkg/rendering"
	"strings"
)

type tvFrame struct{ w, h int }

func newTVFrame(w, h int) rendering.LayerContent { 
	return NewLayerTVFrame(w, h)
}

func (t *tvFrame) GetLayer() rendering.Layer { return rendering.LayerTVFrame }
func (t *tvFrame) GetZ() int                 { return rendering.ZFrame }

func (t *tvFrame) ToLipglossLayer() *lipgloss.Layer {
	// Create TV frame using lipgloss borders instead of manual glyphs
	pad := 3
	innerWidth := t.w - (pad * 2)
	innerHeight := t.h - (pad * 2)
	
	if innerWidth <= 0 || innerHeight <= 0 {
		// Too small for frame, just return empty
		return lipgloss.NewLayer("").X(0).Y(0).Z(t.GetZ()).ID("tv-frame")
	}
	
	// Create frame content
	content := strings.Repeat(" ", innerWidth)
	for i := 1; i < innerHeight; i++ {
		content += "\n" + strings.Repeat(" ", innerWidth)
	}
	
	// Apply thick border style
	frameStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#181c1c")).
		Background(lipgloss.Color("#181c1c")).
		Width(innerWidth).
		Height(innerHeight)
	
	styledContent := frameStyle.Render(content)
	
	return lipgloss.NewLayer(styledContent).
		X(0).
		Y(0).
		Z(t.GetZ()).
		ID("tv-frame")
}
