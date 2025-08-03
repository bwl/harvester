package ui

import "harvester/pkg/rendering"

type testPattern struct{ w,h int }
func newTestPattern(w,h int) *testPattern { return &testPattern{w:w,h:h} }
func (t *testPattern) GetLayer() rendering.Layer { return rendering.LayerGame }
func (t *testPattern) GetZ() int { return rendering.ZPattern }
func (t *testPattern) GetPosition() rendering.Position { return rendering.Position{Horizontal: rendering.Left, Vertical: rendering.Top} }
func (t *testPattern) GetBounds() rendering.Bounds { return rendering.Bounds{Width: t.w, Height: t.h} }
func (t *testPattern) GetGlyphs() [][]rendering.Glyph {
	g := make([][]rendering.Glyph, t.h)
	for y:=0;y<t.h;y++{
		row := make([]rendering.Glyph, t.w)
		for x:=0;x<t.w;x++{
			ch := '.'
			if (x+y)%2==0 { ch=':' }
			row[x] = rendering.Glyph{Char: ch}
		}
		g[y]=row
	}
	return g
}
