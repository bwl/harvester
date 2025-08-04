package ui

import (
	"github.com/charmbracelet/lipgloss/v2"
	"harvester/pkg/rendering"
	"strings"
)

type expanseContent struct {
	g    [][]rendering.Glyph
	w, h int
}

func NewExpanseContent(g [][]rendering.Glyph) *expanseContent {
	if g == nil {
		return nil
	}
	return &expanseContent{g: g, w: len(g[0]), h: len(g)}
}

func (t *expanseContent) GetLayer() rendering.Layer { return rendering.LayerGame }
func (t *expanseContent) GetZ() int                 { return rendering.ZContent }

func (t *expanseContent) ToLipglossLayer() *lipgloss.Layer {
	if t.g == nil || len(t.g) == 0 {
		return lipgloss.NewLayer("").X(0).Y(0).Z(t.GetZ()).ID("expanse-content")
	}

	var content strings.Builder
	for y := 0; y < len(t.g); y++ {
		if y > 0 {
			content.WriteString("\n")
		}
		for x := 0; x < len(t.g[y]); x++ {
			glyph := t.g[y][x]
			
			// Convert glyph color to hex string
			fgHex := "#" + ColorToHex(glyph.Foreground)
			bgHex := "#" + ColorToHex(glyph.Background)
			
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
		Z(t.GetZ()).
		ID("expanse-content")
}

