package ui

import "harvester/pkg/rendering"

type tvFrame struct{ w, h int }

func newTVFrame(w, h int) *tvFrame           { return &tvFrame{w: w, h: h} }
func (t *tvFrame) GetLayer() rendering.Layer { return rendering.LayerTVFrame }
func (t *tvFrame) GetZ() int                 { return rendering.ZFrame }
func (t *tvFrame) GetPosition() rendering.Position {
	return rendering.Position{Horizontal: rendering.Left, Vertical: rendering.Top}
}
func (t *tvFrame) GetBounds() rendering.Bounds { return rendering.Bounds{Width: t.w, Height: t.h} }
func (t *tvFrame) GetGlyphs() [][]rendering.Glyph {
	g := make([][]rendering.Glyph, t.h)
	for y := 0; y < t.h; y++ {
		g[y] = make([]rendering.Glyph, t.w)
	}
	black := rendering.Color{R: 0, G: 0, B: 0}
	pad := 3
	// top and bottom borders
	for y := 0; y < pad && y < t.h; y++ {
		for x := 0; x < t.w; x++ {
			g[y][x] = rendering.Glyph{Char: '█', Foreground: black, Background: black, Alpha: 1.0}
		}
	}
	for y := t.h - pad; y < t.h; y++ {
		if y < 0 {
			continue
		}
		for x := 0; x < t.w; x++ {
			g[y][x] = rendering.Glyph{Char: '█', Foreground: black, Background: black, Alpha: 1.0}
		}
	}
	// left and right borders
	for y := pad; y < t.h-pad; y++ {
		for x := 0; x < pad && x < t.w; x++ {
			g[y][x] = rendering.Glyph{Char: '█', Foreground: black, Background: black, Alpha: 1.0}
		}
		for x := t.w - pad; x < t.w; x++ {
			if x >= 0 {
				g[y][x] = rendering.Glyph{Char: '█', Foreground: black, Background: black, Alpha: 1.0}
			}
		}
	}
	return g
}
