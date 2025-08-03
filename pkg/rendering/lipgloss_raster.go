package rendering

// Parse a subset of ANSI SGR to build glyphs with style and rgb colors.
// Supports: bold(1), italic(3), underline(4), dim(2), reverse(7), reset(0)
// Colors: 38;2;r;g;b (fg truecolor), 48;2;r;g;b (bg truecolor)
func RenderLipglossString(lines []string, defaultFG Color, defaultBG Color, defaultStyle Style) [][]Glyph {
	maxW := 0
	for _, s := range lines { if lw := len([]rune(s)); lw > maxW { maxW = lw } }
	m := make([][]Glyph, len(lines))
	for y, s := range lines {
		row := make([]Glyph, maxW)
		fg, bg, st := defaultFG, defaultBG, defaultStyle
		in := []rune(s)
		i := 0
		x := 0
		for i < len(in) && x < maxW {
			if in[i] == 0x1b && i+1 < len(in) && in[i+1] == '[' {
				// parse SGR
				i += 2
				params := []int{}
				v := 0
				acc := false
				for i < len(in) {
					r := in[i]
					if r >= '0' && r <= '9' { v = v*10 + int(r-'0'); acc = true; i++; continue }
					if r == ';' { params = append(params, v); v = 0; acc = false; i++; continue }
					if r == 'm' { if acc || len(params)==0 { params = append(params, v) }; i++; break }
					// unknown, break
					i++
				}
				if len(params) == 0 { continue }
				// handle params
				for p := 0; p < len(params); p++ {
					sw := params[p]
					switch sw {
					case 0:
						fg, bg, st = defaultFG, defaultBG, StyleNone
					case 1:
						st |= StyleBold
					case 2:
						st |= StyleDim
					case 3:
						st |= StyleItalic
					case 4:
						st |= StyleUnderline
					case 7:
						st |= StyleReverse
					case 38:
						if p+4 < len(params) && params[p+1] == 2 {
							fg = Color{R:uint8(params[p+2]), G:uint8(params[p+3]), B:uint8(params[p+4])}
							p += 4
						}
					case 48:
						if p+4 < len(params) && params[p+1] == 2 {
							bg = Color{R:uint8(params[p+2]), G:uint8(params[p+3]), B:uint8(params[p+4])}
							p += 4
						}
					}
				}
				continue
			}
			ch := in[i]
			row[x] = Glyph{Char: ch, Foreground: fg, Background: bg, Style: st}
			i++
			x++
		}
		// pad remainder
		for ; x < maxW; x++ { row[x] = Glyph{Char: ' ', Foreground: fg, Background: bg, Style: st} }
		m[y] = row
	}
	return m
}
