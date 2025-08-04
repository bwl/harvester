package rendering

// Rendering layers for organizing content by type
type Layer int

const (
	LayerGame Layer = iota
	LayerUI
	LayerHUD
	LayerMenu
	LayerTVFrame
)

// Rect represents a rectangular area (for compatibility)
type Rect struct {
	X, Y, W, H int
}

// GlyphMatrix represents a 2D matrix (deprecated, kept for compatibility)
type GlyphMatrix struct {
	W, H int
	Data [][]interface{} // Empty interface since we don't use glyphs anymore
}

func NewGlyphMatrix(w, h int) *GlyphMatrix {
	return &GlyphMatrix{
		W: w, H: h,
		Data: make([][]interface{}, h),
	}
}

func (g *GlyphMatrix) Clear() {
	// No-op since we don't use glyph data
}

func (g *GlyphMatrix) InBounds(x, y int) bool {
	return x >= 0 && x < g.W && y >= 0 && y < g.H
}

func (g *GlyphMatrix) SetGlyph(x, y int, glyph interface{}) {
	// No-op since we don't use glyph data
}

func (g *GlyphMatrix) GetGlyph(x, y int) (interface{}, bool) {
	if g.InBounds(x, y) {
		return nil, true
	}
	return nil, false
}
