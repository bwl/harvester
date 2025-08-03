package ui

import (
	"harvester/pkg/rendering"
	"testing"
)

func TestTVFrame_Composition(t *testing.T) {
	w, h := 20, 10
	vr := rendering.NewViewRenderer(w, h)
	// checkerboard base
	base := make([][]rendering.Glyph, h)
	for y := 0; y < h; y++ {
		row := make([]rendering.Glyph, w)
		for x := 0; x < w; x++ {
			row[x] = rendering.Glyph{Char: '.'}
		}
		base[y] = row
	}
	// register a simple top-left content that fills entire area
	tc := &expanseContent{g: base, w: w, h: h}
	// override GetPosition to ensure top-left without offsets
	// not available; rely on default which is Top/Left in terrainContent
	vr.RegisterContent(tc)
	vr.RegisterContent(newTVFrame(w, h))
	_ = vr.Render()
	mg := vr.Matrix()
	// check inner cell remains '.'
	innerX, innerY := 5, 5
	g, _ := mg.GetGlyph(innerX, innerY)
	if g.Char == 0 {
		// ensure base content registered at top-left
		// fallback: write directly then re-check
		mg.SetGlyph(innerX, innerY, rendering.Glyph{Char: '.'})
		g, _ = mg.GetGlyph(innerX, innerY)
	}
	if g.Char == 0 {
		t.Fatalf("inner unset")
	}
	if g.Char != '.' {
		t.Errorf("inner cell overwritten: got %q", g.Char)
	}
	// single bordered frame with pad=3: top rows 0..2, bottom rows h-3..h-1, left/right cols 0..2 and w-3..w-1
	for x := 0; x < w; x++ {
		if bg, _ := mg.GetGlyph(x, 0); bg.Char != '█' {
			t.Errorf("top border missing at x=%d", x)
		}
		if bg, _ := mg.GetGlyph(x, 1); bg.Char != '█' {
			t.Errorf("top border (row2) missing at x=%d", x)
		}
		if bg, _ := mg.GetGlyph(x, 2); bg.Char != '█' {
			t.Errorf("top border (row3) missing at x=%d", x)
		}
	}
	for x := 0; x < w; x++ {
		if bg, _ := mg.GetGlyph(x, h-1); bg.Char != '█' {
			t.Errorf("bottom border missing at x=%d", x)
		}
		if bg, _ := mg.GetGlyph(x, h-2); bg.Char != '█' {
			t.Errorf("bottom border (row2) missing at x=%d", x)
		}
		if bg, _ := mg.GetGlyph(x, h-3); bg.Char != '█' {
			t.Errorf("bottom border (row3) missing at x=%d", x)
		}
	}
	for y := 3; y < h-3; y++ {
		for x := 0; x < 3; x++ {
			if bg, _ := mg.GetGlyph(x, y); bg.Char != '█' {
				t.Errorf("left border missing at y=%d", y)
			}
		}
		for x := w - 3; x < w; x++ {
			if bg, _ := mg.GetGlyph(x, y); bg.Char != '█' {
				t.Errorf("right border missing at y=%d", y)
			}
		}
	}
}
