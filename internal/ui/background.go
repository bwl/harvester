package ui

import (
	"math/rand"
	"harvester/pkg/rendering"
)

type backgroundLayer struct{ w,h int; rng *rand.Rand }
func newBackgroundLayer(w,h int) *backgroundLayer { return &backgroundLayer{w:w,h:h, rng: rand.New(rand.NewSource(42))} }
func (b *backgroundLayer) GetLayer() rendering.Layer { return rendering.LayerGame }
func (b *backgroundLayer) GetZ() int { return rendering.ZBackground }
func (b *backgroundLayer) GetPosition() rendering.Position { return rendering.Position{Horizontal: rendering.Left, Vertical: rendering.Top} }
func (b *backgroundLayer) GetBounds() rendering.Bounds { return rendering.Bounds{Width: b.w, Height: b.h} }
func (b *backgroundLayer) GetGlyphs() [][]rendering.Glyph {
	g := make([][]rendering.Glyph, b.h)
	for y:=0;y<b.h;y++{
		row := make([]rendering.Glyph, b.w)
		for x:=0;x<b.w;x++{
			v := uint8( b.rng.Intn(32) ) // 0..31 dark grey
			c := rendering.Color{R:v,G:v,B:v}
			row[x] = rendering.Glyph{Char:'█', Foreground:c, Background:c}
		}
		g[y]=row
	}
	return g
}
