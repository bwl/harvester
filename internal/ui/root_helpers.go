package ui

import "harvester/pkg/rendering"

type textBlock struct{ g [][]rendering.Glyph; w,h int }
func newTextBlock(s string, w,h int) *textBlock {
	lines := splitLines(s)
	g := rendering.RenderLipglossString(lines, rendering.Color{}, rendering.Color{}, rendering.StyleNone)
	return &textBlock{g:g,w:w,h:h}
}
func (t *textBlock) GetLayer() rendering.Layer { return rendering.LayerMenu }
func (t *textBlock) GetZ() int { return rendering.ZContent + 10 }
func (t *textBlock) GetPosition() rendering.Position { return rendering.Position{Horizontal: rendering.Left, Vertical: rendering.Top} }
func (t *textBlock) GetBounds() rendering.Bounds { return rendering.Bounds{Width: t.w, Height: t.h} }
func (t *textBlock) GetGlyphs() [][]rendering.Glyph { return t.g }
