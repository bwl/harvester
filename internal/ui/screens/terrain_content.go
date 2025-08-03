package screens

import "harvester/pkg/rendering"

type TerrainContent struct {
	glyphs [][]rendering.Glyph
	w, h   int
}

func NewTerrainContent(g [][]rendering.Glyph, w, h int) *TerrainContent {
	return &TerrainContent{glyphs: g, w: w, h: h}
}

func (t *TerrainContent) GetLayer() rendering.Layer { return rendering.LayerGame }
func (t *TerrainContent) GetZ() int { return rendering.ZContent }
func (t *TerrainContent) GetPosition() rendering.Position {
	return rendering.Position{Horizontal: rendering.Left, Vertical: rendering.Top}
}
func (t *TerrainContent) GetBounds() rendering.Bounds {
	return rendering.Bounds{Width: t.w, Height: t.h}
}
func (t *TerrainContent) GetGlyphs() [][]rendering.Glyph { return t.glyphs }
