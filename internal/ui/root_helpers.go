package ui

import (
	"github.com/charmbracelet/lipgloss/v2"
	"harvester/pkg/rendering"
)

type textBlock struct {
	text string
	w, h int
}

func newTextBlock(s string, w, h int) *textBlock {
	return &textBlock{text: s, w: w, h: h}
}

func (t *textBlock) GetLayer() rendering.Layer { return rendering.LayerMenu }
func (t *textBlock) GetZ() int                 { return rendering.ZContent + 10 }

func (t *textBlock) ToLipglossLayer() *lipgloss.Layer {
	style := lipgloss.NewStyle().
		Width(t.w).
		Height(t.h).
		Foreground(lipgloss.Color("#ffffff"))

	styledText := style.Render(t.text)

	return lipgloss.NewLayer(styledText).
		X(0).
		Y(0).
		Z(t.GetZ()).
		ID("text-block")
}

