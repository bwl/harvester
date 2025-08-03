package ui

import "harvester/pkg/rendering"

type expanseContent struct {
	g    [][]rendering.Glyph
	w, h int
}

func newExpanseContent(g [][]rendering.Glyph) *expanseContent {
	if g == nil {
		return nil
	}
	return &expanseContent{g: g, w: len(g[0]), h: len(g)}
}
func (t *expanseContent) GetLayer() rendering.Layer { return rendering.LayerGame }
func (t *expanseContent) GetZ() int                 { return rendering.ZContent }
func (t *expanseContent) GetPosition() rendering.Position {
	return rendering.Position{Horizontal: rendering.Left, Vertical: rendering.Top}
}
func (t *expanseContent) GetBounds() rendering.Bounds {
	return rendering.Bounds{Width: t.w, Height: t.h}
}
func (t *expanseContent) GetGlyphs() [][]rendering.Glyph { return t.g }
