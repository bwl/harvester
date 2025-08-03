package ui

import "harvester/pkg/rendering"

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
func (t *planetSurfaceContent) GetPosition() rendering.Position {
	return rendering.Position{Horizontal: rendering.Left, Vertical: rendering.Top}
}
func (t *planetSurfaceContent) GetBounds() rendering.Bounds {
	return rendering.Bounds{Width: t.w, Height: t.h}
}
func (t *planetSurfaceContent) GetGlyphs() [][]rendering.Glyph { return t.g }
