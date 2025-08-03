package rendering

import (
	"github.com/charmbracelet/lipgloss"
	"harvester/pkg/systems"
)

// DrawableToGlyph converts a systems.Drawable to a rendering.Glyph
func DrawableToGlyph(d systems.Drawable) Glyph {
	return Glyph{
		Char:       d.Glyph,
		Foreground: lipglossColorToRenderingColor(d.Style.GetForeground()),
		Background: lipglossColorToRenderingColor(d.Style.GetBackground()),
		Style:      lipglossStyleToRenderingStyle(d.Style),
		Alpha:      d.Alpha,
		BlendMode:  d.BlendMode,
	}
}

// lipglossColorToRenderingColor converts lipgloss color to rendering color
func lipglossColorToRenderingColor(lc lipgloss.TerminalColor) Color {
	if lc == nil {
		return Color{} // Empty color
	}

	// Handle different lipgloss color types
	switch color := lc.(type) {
	case lipgloss.Color:
		// Parse hex color or ANSI color
		return parseColorString(string(color))
	case lipgloss.AdaptiveColor:
		// Use light color for now (could be made context-aware)
		return parseColorString(string(color.Light))
	default:
		return Color{} // Unknown color type
	}
}

// parseColorString parses color strings (hex, ANSI codes, etc.)
func parseColorString(colorStr string) Color {
	// Handle hex colors
	if len(colorStr) == 7 && colorStr[0] == '#' {
		return parseHexColor(colorStr)
	}

	// Handle ANSI color codes
	if ansiColor, ok := ansiToRGB(colorStr); ok {
		return ansiColor
	}

	// Default to empty color
	return Color{}
}

// parseHexColor parses hex color strings like "#FF0000"
func parseHexColor(hex string) Color {
	if len(hex) != 7 || hex[0] != '#' {
		return Color{}
	}

	// Simple hex parsing
	var r, g, b uint8

	// Parse each component
	if rVal, ok := parseHexByte(hex[1:3]); ok {
		r = rVal
	}
	if gVal, ok := parseHexByte(hex[3:5]); ok {
		g = gVal
	}
	if bVal, ok := parseHexByte(hex[5:7]); ok {
		b = bVal
	}

	return Color{R: r, G: g, B: b}
}

// parseHexByte parses a 2-character hex string to byte
func parseHexByte(hex string) (uint8, bool) {
	if len(hex) != 2 {
		return 0, false
	}

	var result uint8
	for i, c := range hex {
		var val uint8
		switch {
		case c >= '0' && c <= '9':
			val = uint8(c - '0')
		case c >= 'A' && c <= 'F':
			val = uint8(c - 'A' + 10)
		case c >= 'a' && c <= 'f':
			val = uint8(c - 'a' + 10)
		default:
			return 0, false
		}

		if i == 0 {
			result = val * 16
		} else {
			result += val
		}
	}

	return result, true
}

// ansiToRGB converts ANSI color codes to RGB
func ansiToRGB(ansi string) (Color, bool) {
	// Basic ANSI colors
	ansiColors := map[string]Color{
		"0":  {R: 0, G: 0, B: 0},       // Black
		"1":  {R: 128, G: 0, B: 0},     // Dark Red
		"2":  {R: 0, G: 128, B: 0},     // Dark Green
		"3":  {R: 128, G: 128, B: 0},   // Dark Yellow
		"4":  {R: 0, G: 0, B: 128},     // Dark Blue
		"5":  {R: 128, G: 0, B: 128},   // Dark Magenta
		"6":  {R: 0, G: 128, B: 128},   // Dark Cyan
		"7":  {R: 192, G: 192, B: 192}, // Light Gray
		"8":  {R: 128, G: 128, B: 128}, // Dark Gray
		"9":  {R: 255, G: 0, B: 0},     // Red
		"10": {R: 0, G: 255, B: 0},     // Green
		"11": {R: 255, G: 255, B: 0},   // Yellow
		"12": {R: 0, G: 0, B: 255},     // Blue
		"13": {R: 255, G: 0, B: 255},   // Magenta
		"14": {R: 0, G: 255, B: 255},   // Cyan
		"15": {R: 255, G: 255, B: 255}, // White

		// Extended colors (some common ones)
		"226": {R: 255, G: 255, B: 0},   // Bright Yellow
		"27":  {R: 0, G: 135, B: 175},   // Blue
		"240": {R: 88, G: 88, B: 88},    // Dark Gray
		"244": {R: 128, G: 128, B: 128}, // Gray
		"252": {R: 208, G: 208, B: 208}, // Light Gray
		"235": {R: 38, G: 38, B: 38},    // Very Dark Gray
	}

	if color, exists := ansiColors[ansi]; exists {
		return color, true
	}

	return Color{}, false
}

// lipglossStyleToRenderingStyle converts lipgloss style attributes to rendering style
func lipglossStyleToRenderingStyle(style lipgloss.Style) Style {
	var result Style = StyleNone

	// Extract style attributes from lipgloss style
	// This is a simplified conversion - lipgloss doesn't expose style flags directly
	// In practice, we might need to inspect the style more carefully

	// For now, return basic style
	// TODO: Implement proper style extraction from lipgloss.Style
	return result
}

// DrawablesToGlyphs converts a slice of Drawables to a 2D glyph matrix
func DrawablesToGlyphs(drawables []systems.Drawable, width, height int) [][]Glyph {
	matrix := make([][]Glyph, height)
	for i := range matrix {
		matrix[i] = make([]Glyph, width)
		// Initialize with default alpha 1.0 for empty spaces
		for j := range matrix[i] {
			matrix[i][j] = Glyph{Alpha: 1.0}
		}
	}

	// Place drawables in matrix
	for _, drawable := range drawables {
		x, y := drawable.X, drawable.Y
		if x >= 0 && x < width && y >= 0 && y < height {
			matrix[y][x] = DrawableToGlyph(drawable)
		}
	}

	return matrix
}
