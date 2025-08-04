package screens

import (
	"github.com/charmbracelet/lipgloss/v2"
	"harvester/pkg/rendering"
	"strings"
)

type TerrainContent struct {
	glyphs [][]rendering.Glyph
	w, h   int
}

func NewTerrainContent(g [][]rendering.Glyph, w, h int) *TerrainContent {
	return &TerrainContent{glyphs: g, w: w, h: h}
}

func (t *TerrainContent) GetLayer() rendering.Layer { return rendering.LayerGame }
func (t *TerrainContent) GetZ() int                 { return rendering.ZContent }
func (t *TerrainContent) GetPosition() rendering.Position {
	return rendering.Position{Horizontal: rendering.Left, Vertical: rendering.Top}
}
func (t *TerrainContent) GetBounds() rendering.Bounds {
	return rendering.Bounds{Width: t.w, Height: t.h}
}
func (t *TerrainContent) GetGlyphs() [][]rendering.Glyph { return t.glyphs }

func (t *TerrainContent) ToLipglossLayer() *lipgloss.Layer {
	if t.glyphs == nil || len(t.glyphs) == 0 {
		return lipgloss.NewLayer("").X(0).Y(0).Z(t.GetZ()).ID("terrain")
	}

	var content strings.Builder
	for y := 0; y < len(t.glyphs); y++ {
		if y > 0 {
			content.WriteString("\n")
		}
		for x := 0; x < len(t.glyphs[y]); x++ {
			glyph := t.glyphs[y][x]
			
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
		Z(t.GetZ()).
		ID("terrain")
}


