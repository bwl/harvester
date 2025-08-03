package rendering

import (
	"testing"
)

type dummyContent struct{ w,h int }
func (d dummyContent) GetLayer() Layer { return LayerGame }
func (d dummyContent) GetZ() int { return 0 }
func (d dummyContent) GetPosition() Position { return Position{Horizontal: Left, Vertical: Top} }
func (d dummyContent) GetBounds() Bounds { return Bounds{Width: d.w, Height: d.h} }
func (d dummyContent) GetGlyphs() [][]Glyph {
	g := make([][]Glyph, d.h)
	for y := 0; y < d.h; y++ {
		row := make([]Glyph, d.w)
		for x := 0; x < d.w; x++ { row[x] = Glyph{Char: '.'} }
		g[y] = row
	}
	return g
}

func BenchmarkRender_80x24(b *testing.B) {
	vr := NewViewRenderer(80,24)
	vr.RegisterContent(dummyContent{w:80,h:24})
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		_ = vr.Render()
	}
}

func BenchmarkRender_120x40(b *testing.B) {
	vr := NewViewRenderer(120,40)
	vr.RegisterContent(dummyContent{w:120,h:40})
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		_ = vr.Render()
	}
}
