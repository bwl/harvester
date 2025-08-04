package screens

import (
	"github.com/charmbracelet/lipgloss/v2"
	"harvester/pkg/rendering"
	"strings"
)

type MenuContent struct {
	glyphs [][]rendering.Glyph
	w, h   int
}

func NewMenuContent(g [][]rendering.Glyph, w, h int) *MenuContent {
	return &MenuContent{glyphs: g, w: w, h: h}
}

func (m *MenuContent) GetLayer() rendering.Layer { return rendering.LayerMenu }
func (m *MenuContent) GetZ() int                 { return rendering.ZMenu }
func (m *MenuContent) GetPosition() rendering.Position {
	return rendering.Position{Horizontal: rendering.CenterH, Vertical: rendering.CenterV}
}
func (m *MenuContent) GetBounds() rendering.Bounds    { return rendering.Bounds{Width: m.w, Height: m.h} }
func (m *MenuContent) GetGlyphs() [][]rendering.Glyph { return m.glyphs }

func (m *MenuContent) ToLipglossLayer() *lipgloss.Layer {
	if m.glyphs == nil || len(m.glyphs) == 0 {
		return lipgloss.NewLayer("").X(0).Y(0).Z(m.GetZ()).ID("menu")
	}

	var content strings.Builder
	for y := 0; y < len(m.glyphs); y++ {
		if y > 0 {
			content.WriteString("\n")
		}
		for x := 0; x < len(m.glyphs[y]); x++ {
			glyph := m.glyphs[y][x]
			
			// Convert glyph color to hex string
			fgHex := "#" + colorToHex(glyph.Foreground)
			bgHex := "#" + colorToHex(glyph.Background)
			
			// Create styled character
			styledChar := lipgloss.NewStyle().
				Foreground(lipgloss.Color(fgHex)).
				Background(lipgloss.Color(bgHex)).
				Render(string(glyph.Char))
			
			content.WriteString(styledChar)
		}
	}

	return lipgloss.NewLayer(content.String()).
		X(0).
		Y(0).
		Z(m.GetZ()).
		ID("menu")
}


