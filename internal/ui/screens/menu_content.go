package screens

import "harvester/pkg/rendering"

type MenuContent struct {
	glyphs [][]rendering.Glyph
	w, h   int
}

func NewMenuContent(g [][]rendering.Glyph, w, h int) *MenuContent {
	return &MenuContent{glyphs: g, w: w, h: h}
}

func (m *MenuContent) GetLayer() rendering.Layer { return rendering.LayerMenu }
func (m *MenuContent) GetZ() int { return rendering.ZMenu }
func (m *MenuContent) GetPosition() rendering.Position {
	return rendering.Position{Horizontal: rendering.CenterH, Vertical: rendering.CenterV}
}
func (m *MenuContent) GetBounds() rendering.Bounds    { return rendering.Bounds{Width: m.w, Height: m.h} }
func (m *MenuContent) GetGlyphs() [][]rendering.Glyph { return m.glyphs }
