package rendering

import (
	"fmt"
	"os"
	"time"
	"harvester/pkg/components"
)

type ViewRenderer struct {
	width, height int
	layers        map[Layer][]RenderableContent
	matrix        *GlyphMatrix
	dirty         []Rect
	dirtyAll      bool
}

type LinePatch struct { Y int; Line string }

func (v *ViewRenderer) Matrix() *GlyphMatrix { return v.matrix }

func NewViewRenderer(w, h int) *ViewRenderer {
	vr := &ViewRenderer{
		width:  w,
		height: h,
		layers: make(map[Layer][]RenderableContent),
		matrix: NewGlyphMatrix(w, h),
		dirtyAll: true,
	}
	return vr
}

func (v *ViewRenderer) SetDimensions(w, h int) {
	v.width, v.height = w, h
	v.matrix = NewGlyphMatrix(w, h)
	v.dirty = nil
	v.dirtyAll = true
}

func (v *ViewRenderer) RegisterContent(c RenderableContent) {
	l := c.GetLayer()
	v.layers[l] = append(v.layers[l], c)
}

func (v *ViewRenderer) UnregisterAll() {
	for k := range v.layers { v.layers[k] = v.layers[k][:0] }
	v.MarkDirtyAll()
}

func (v *ViewRenderer) MarkDirty(x,y,w,h int) { v.dirty = append(v.dirty, Rect{X:x,Y:y,W:w,H:h}) }
func (v *ViewRenderer) MarkDirtyAll() { v.dirtyAll = true }
func (v *ViewRenderer) DirtyRegions() []Rect { if v.dirtyAll { return []Rect{{0,0,v.width,v.height}} }; return v.dirty }

func (v *ViewRenderer) Render() string {
	var start time.Time
	if os.Getenv("VR_PROFILE") == "1" { start = time.Now() }
	v.matrix.Clear()
	var all []RenderableContent
	for _, slice := range v.layers { all = append(all, slice...) }
	// stable sort by Z (ascending); equal Z preserve registration order
	for i := 1; i < len(all); i++ {
		j := i
		for j > 0 && all[j-1].GetZ() > all[j].GetZ() {
			all[j-1], all[j] = all[j], all[j-1]
			j--
		}
	}
	for _, c := range all { v.composite(c) }
	s := v.matrixToString()
	v.dirty = nil; v.dirtyAll = false
	if !start.IsZero() {
		fmt.Fprintf(os.Stderr, "[vr] render=%s\n", time.Since(start))
	}
	return s
}

func (v *ViewRenderer) composite(c RenderableContent) {
	b := c.GetBounds()
	p := c.GetPosition()
	sx, sy := CalculatePosition(p, b, v.width, v.height)
	glyphs := c.GetGlyphs()
	
	for y := 0; y < len(glyphs); y++ {
		row := glyphs[y]
		for x := 0; x < len(row); x++ {
			tx, ty := sx+x, sy+y
			if v.matrix.InBounds(tx, ty) {
				newGlyph := row[x]
				
				// Skip completely transparent glyphs
				if newGlyph.Alpha <= 0.0 {
					continue
				}
				
				// Skip empty glyphs (no content to render)
				if newGlyph.Char == 0 && (newGlyph.Foreground == (Color{})) && (newGlyph.Background == (Color{})) && newGlyph.Style == StyleNone {
					continue
				}
				
				// Get existing glyph at position
				existingGlyph, _ := v.matrix.GetGlyph(tx, ty)
				
				// Composite based on alpha
				finalGlyph := v.blendGlyphs(existingGlyph, newGlyph)
				
				v.matrix.SetGlyph(tx, ty, finalGlyph)
				v.MarkDirty(tx, ty, 1, 1)
			}
		}
	}
}

// blendGlyphs performs alpha blending between bottom (existing) and top (new) glyphs
func (v *ViewRenderer) blendGlyphs(bottom, top Glyph) Glyph {
	alpha := top.Alpha
	
	// Fast path: fully opaque - just replace
	if alpha >= 1.0 {
		// Ensure result has alpha 1.0
		result := top
		result.Alpha = 1.0
		return result
	}
	
	// Fast path: fully transparent - keep bottom
	if alpha <= 0.0 {
		return bottom
	}
	
	// Alpha blending required
	result := Glyph{
		Alpha: 1.0, // Result is always opaque after blending
	}
	
	// Character selection strategy
	result.Char = v.selectChar(bottom, top, alpha)
	
	// Color blending based on blend mode
	switch top.BlendMode {
	case components.BlendAdditive:
		result.Foreground = v.blendColorsAdditive(bottom.Foreground, top.Foreground, alpha)
		result.Background = v.blendColorsAdditive(bottom.Background, top.Background, alpha)
		
	case components.BlendMultiply:
		result.Foreground = v.blendColorsMultiply(bottom.Foreground, top.Foreground, alpha)
		result.Background = v.blendColorsMultiply(bottom.Background, top.Background, alpha)
		
	case components.BlendScreen:
		result.Foreground = v.blendColorsScreen(bottom.Foreground, top.Foreground, alpha)
		result.Background = v.blendColorsScreen(bottom.Background, top.Background, alpha)
		
	default: // BlendNormal fallback
		result.Foreground = v.blendColorsNormal(bottom.Foreground, top.Foreground, alpha)
		result.Background = v.blendColorsNormal(bottom.Background, top.Background, alpha)
	}
	
	// Style blending
	result.Style = v.blendStyles(bottom.Style, top.Style, alpha)
	
	return result
}

// selectChar chooses which character to use based on alpha
func (v *ViewRenderer) selectChar(bottom, top Glyph, alpha float64) rune {
	// Threshold-based selection
	if alpha > 0.5 {
		if top.Char != 0 {
			return top.Char
		}
	}
	
	if bottom.Char != 0 {
		return bottom.Char
	}
	
	return ' ' // Default fallback
}

// blendColorsNormal performs standard alpha blending
func (v *ViewRenderer) blendColorsNormal(bottom, top Color, alpha float64) Color {
	if (top == Color{}) {
		return bottom // No color to blend
	}
	if (bottom == Color{}) {
		return Color{
			R: uint8(float64(top.R) * alpha),
			G: uint8(float64(top.G) * alpha),
			B: uint8(float64(top.B) * alpha),
		}
	}
	
	// Standard alpha blending: result = bottom*(1-alpha) + top*alpha
	return Color{
		R: uint8(float64(bottom.R)*(1-alpha) + float64(top.R)*alpha),
		G: uint8(float64(bottom.G)*(1-alpha) + float64(top.G)*alpha),
		B: uint8(float64(bottom.B)*(1-alpha) + float64(top.B)*alpha),
	}
}

// blendColorsAdditive performs additive blending
func (v *ViewRenderer) blendColorsAdditive(bottom, top Color, alpha float64) Color {
	if (top == Color{}) {
		return bottom
	}
	
	// Additive: result = bottom + top*alpha (clamped to 255)
	return Color{
		R: uint8(min(255, int(bottom.R)+int(float64(top.R)*alpha))),
		G: uint8(min(255, int(bottom.G)+int(float64(top.G)*alpha))),
		B: uint8(min(255, int(bottom.B)+int(float64(top.B)*alpha))),
	}
}

// blendColorsMultiply performs multiply blending
func (v *ViewRenderer) blendColorsMultiply(bottom, top Color, alpha float64) Color {
	if (top == Color{}) {
		return bottom
	}
	
	// Multiply: result = bottom * (top*alpha/255)
	return Color{
		R: uint8(float64(bottom.R) * (float64(top.R)*alpha/255.0)),
		G: uint8(float64(bottom.G) * (float64(top.G)*alpha/255.0)),
		B: uint8(float64(bottom.B) * (float64(top.B)*alpha/255.0)),
	}
}

// blendColorsScreen performs screen blending
func (v *ViewRenderer) blendColorsScreen(bottom, top Color, alpha float64) Color {
	if (top == Color{}) {
		return bottom
	}
	
	// Screen: result = 255 - (255-bottom) * (255-top*alpha) / 255
	return Color{
		R: uint8(255 - (255-int(bottom.R))*(255-int(float64(top.R)*alpha))/255),
		G: uint8(255 - (255-int(bottom.G))*(255-int(float64(top.G)*alpha))/255),
		B: uint8(255 - (255-int(bottom.B))*(255-int(float64(top.B)*alpha))/255),
	}
}

// blendStyles combines styles based on alpha threshold
func (v *ViewRenderer) blendStyles(bottom, top Style, alpha float64) Style {
	if alpha > 0.5 {
		// Above threshold: combine styles
		return bottom | top
	} else {
		// Below threshold: keep bottom style
		return bottom
	}
}

func (v *ViewRenderer) matrixToString() string {
	var b []byte
	reset := "\x1b[0m"
	current := struct {
		fg   Color
		bg   Color
		bold bool
	}{}
	emitStyle := func(g Glyph) {
		seq := ""
		bold := g.Style&StyleBold != 0
		if bold != current.bold {
			if bold { seq += "\x1b[1m" } else { seq += reset }
			current.bold = bold
		}
		if g.Foreground != current.fg {
			if (g.Foreground == Color{}) { seq += reset } else { seq += fmt.Sprintf("\x1b[38;2;%d;%d;%dm", g.Foreground.R, g.Foreground.G, g.Foreground.B) }
			current.fg = g.Foreground
		}
		if g.Background != current.bg {
			if (g.Background == Color{}) { seq += reset } else { seq += fmt.Sprintf("\x1b[48;2;%d;%d;%dm", g.Background.R, g.Background.G, g.Background.B) }
			current.bg = g.Background
		}
		if seq != "" { b = append(b, []byte(seq)...)}
	}
	var start time.Time
	if os.Getenv("VR_PROFILE") == "1" { start = time.Now() }
	for y := 0; y < v.height; y++ {
		// reset style at line start
		b = append(b, []byte(reset)...)
		current = struct {
			fg   Color
			bg   Color
			bold bool
		}{}
		for x := 0; x < v.width; x++ {
			g, _ := v.matrix.GetGlyph(x, y)
			w := runeWidth(g.Char)
			if w == 0 {
				// combining mark: append to previous cell output if any
				b = append(b, string(g.Char)...)
				continue
			}
			if w == 2 {
				if g.Char == 0 { g.Char = ' ' }
				emitStyle(g)
				b = append(b, string(g.Char)...)
				x++ // skip next cell
				continue
			}
			if g.Char == 0 { g.Char = ' ' }
			emitStyle(g)
			b = append(b, string(g.Char)...)
		}
		b = append(b, []byte(reset)...)
		b = append(b, '\n')
	}
	if !start.IsZero() { fmt.Fprintf(os.Stderr, "[vr] stringify=%s\n", time.Since(start)) }
	return string(b)
}

func (v *ViewRenderer) matrixLine(y int) string {
	var b []byte
	reset := "\x1b[0m"
	current := struct{ fg, bg Color; bold bool }{}
	emitStyle := func(g Glyph) {
		seq := ""
		bold := g.Style&StyleBold != 0
		if bold != current.bold { if bold { seq += "\x1b[1m" } else { seq += reset }; current.bold = bold }
		if g.Foreground != current.fg { if (g.Foreground == Color{}) { seq += reset } else { seq += fmt.Sprintf("\x1b[38;2;%d;%d;%dm", g.Foreground.R, g.Foreground.G, g.Foreground.B) }; current.fg = g.Foreground }
		if g.Background != current.bg { if (g.Background == Color{}) { seq += reset } else { seq += fmt.Sprintf("\x1b[48;2;%d;%d;%dm", g.Background.R, g.Background.G, g.Background.B) }; current.bg = g.Background }
		if seq != "" { b = append(b, []byte(seq)...)}
	}
	b = append(b, []byte(reset)...)
	current = struct{ fg, bg Color; bold bool }{}
	for x := 0; x < v.width; x++ {
		g, _ := v.matrix.GetGlyph(x, y)
		w := runeWidth(g.Char)
		if w == 0 { b = append(b, string(g.Char)...); continue }
		if w == 2 {
			if g.Char == 0 { g.Char = ' ' }
			emitStyle(g)
			b = append(b, string(g.Char)...)
			x++
			continue
		}
		if g.Char == 0 { g.Char = ' ' }
		emitStyle(g)
		b = append(b, string(g.Char)...)
	}
	b = append(b, []byte(reset)...)
	return string(b)
}

func (v *ViewRenderer) RenderPatch() []LinePatch {
	if v.dirtyAll { v.dirty = []Rect{{0,0,v.width,v.height}} }
	// collect unique lines
	seen := make(map[int]struct{})
	var lines []int
	for _, r := range v.dirty {
		for y := r.Y; y < r.Y+r.H; y++ {
			if y < 0 || y >= v.height { continue }
			if _, ok := seen[y]; ok { continue }
			seen[y] = struct{}{}
			lines = append(lines, y)
		}
	}
	if len(lines) == 0 { return nil }
	// build patches
	patches := make([]LinePatch, 0, len(lines))
	for _, y := range lines {
		patches = append(patches, LinePatch{Y: y, Line: v.matrixLine(y)})
	}
	v.dirty = nil; v.dirtyAll = false
	return patches
}
