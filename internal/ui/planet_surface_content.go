package ui

import (
	"github.com/charmbracelet/lipgloss/v2"
	"harvester/pkg/rendering"
	"strings"
)

type planetSurfaceContent struct {
	g    [][]rendering.Glyph
	w, h int
}

func newPlanetSurfaceContent(g [][]rendering.Glyph) *planetSurfaceContent {
	if g == nil {
		return nil
	}
	return &planetSurfaceContent{g: g, w: len(g[0]), h: len(g)}
}

func (t *planetSurfaceContent) GetLayer() rendering.Layer { return rendering.LayerGame }
func (t *planetSurfaceContent) GetZ() int                 { return rendering.ZContent }

func (t *planetSurfaceContent) ToLipglossLayer() *lipgloss.Layer {
	if t.g == nil || len(t.g) == 0 {
		return lipgloss.NewLayer("").X(0).Y(0).Z(t.GetZ()).ID("planet-surface")
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
		ID("planet-surface")
}

