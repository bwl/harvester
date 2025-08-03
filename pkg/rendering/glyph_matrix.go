package rendering

type GlyphMatrix struct {
	W, H int
	Data [][]Glyph
}

func NewGlyphMatrix(w, h int) *GlyphMatrix {
	gm := &GlyphMatrix{W: w, H: h, Data: make([][]Glyph, h)}
	for i := range gm.Data {
		gm.Data[i] = make([]Glyph, w)
	}
	return gm
}

func (g *GlyphMatrix) Clear() {
	for y := 0; y < g.H; y++ {
		row := g.Data[y]
		for x := 0; x < g.W; x++ {
			row[x] = Glyph{}
		}
	}
}

func (g *GlyphMatrix) InBounds(x, y int) bool {
	return x >= 0 && x < g.W && y >= 0 && y < g.H
}

func (g *GlyphMatrix) SetGlyph(x, y int, glyph Glyph) {
	if g.InBounds(x, y) {
		g.Data[y][x] = glyph
	}
}

func (g *GlyphMatrix) GetGlyph(x, y int) (Glyph, bool) {
	if g.InBounds(x, y) {
		return g.Data[y][x], true
	}
	return Glyph{}, false
}
