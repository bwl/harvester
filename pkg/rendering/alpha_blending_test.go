package rendering

import (
	"harvester/pkg/components"
	"testing"
)

func TestAlphaBlending(t *testing.T) {
	vr := NewViewRenderer(10, 10)

	// Test color blending normal mode
	bottom := Color{R: 100, G: 100, B: 100}
	top := Color{R: 200, G: 200, B: 200}
	alpha := 0.5

	result := vr.blendColorsNormal(bottom, top, alpha)
	expected := Color{
		R: uint8(100*0.5 + 200*0.5), // 150
		G: uint8(100*0.5 + 200*0.5), // 150
		B: uint8(100*0.5 + 200*0.5), // 150
	}

	if result != expected {
		t.Errorf("Normal blend failed: expected %+v, got %+v", expected, result)
	}
}

func TestAlphaBlendingAdditive(t *testing.T) {
	vr := NewViewRenderer(10, 10)

	// Test additive blending
	bottom := Color{R: 100, G: 100, B: 100}
	top := Color{R: 100, G: 100, B: 100}
	alpha := 0.5

	result := vr.blendColorsAdditive(bottom, top, alpha)
	expected := Color{
		R: uint8(100 + 100*0.5), // 150
		G: uint8(100 + 100*0.5), // 150
		B: uint8(100 + 100*0.5), // 150
	}

	if result != expected {
		t.Errorf("Additive blend failed: expected %+v, got %+v", expected, result)
	}
}

func TestAlphaBlendingClamping(t *testing.T) {
	vr := NewViewRenderer(10, 10)

	// Test clamping in additive mode
	bottom := Color{R: 200, G: 200, B: 200}
	top := Color{R: 200, G: 200, B: 200}
	alpha := 1.0

	result := vr.blendColorsAdditive(bottom, top, alpha)
	expected := Color{R: 255, G: 255, B: 255} // Should clamp to 255

	if result != expected {
		t.Errorf("Additive clamping failed: expected %+v, got %+v", expected, result)
	}
}

func TestGlyphBlending(t *testing.T) {
	vr := NewViewRenderer(10, 10)

	bottom := Glyph{
		Char:       'A',
		Foreground: Color{R: 100, G: 100, B: 100},
		Background: Color{R: 50, G: 50, B: 50},
		Style:      StyleBold,
		Alpha:      1.0,
	}

	top := Glyph{
		Char:       'B',
		Foreground: Color{R: 200, G: 200, B: 200},
		Background: Color{R: 150, G: 150, B: 150},
		Style:      StyleItalic,
		Alpha:      0.5,
		BlendMode:  components.BlendNormal,
	}

	result := vr.blendGlyphs(bottom, top)

	// At 50% alpha, should use bottom character
	if result.Char != 'A' {
		t.Errorf("Expected char 'A', got '%c'", result.Char)
	}

	// Result should always be opaque
	if result.Alpha != 1.0 {
		t.Errorf("Expected alpha 1.0, got %f", result.Alpha)
	}

	// Colors should be blended
	expectedFg := Color{
		R: uint8(100*0.5 + 200*0.5), // 150
		G: uint8(100*0.5 + 200*0.5), // 150
		B: uint8(100*0.5 + 200*0.5), // 150
	}

	if result.Foreground != expectedFg {
		t.Errorf("Expected foreground %+v, got %+v", expectedFg, result.Foreground)
	}
}

func TestGlyphBlendingOpaque(t *testing.T) {
	vr := NewViewRenderer(10, 10)

	bottom := Glyph{
		Char:       'A',
		Foreground: Color{R: 100, G: 100, B: 100},
		Alpha:      1.0,
	}

	top := Glyph{
		Char:       'B',
		Foreground: Color{R: 200, G: 200, B: 200},
		Alpha:      1.0, // Fully opaque
		BlendMode:  components.BlendNormal,
	}

	result := vr.blendGlyphs(bottom, top)

	// Should completely replace bottom
	if result.Char != 'B' {
		t.Errorf("Expected char 'B', got '%c'", result.Char)
	}

	if result.Foreground != top.Foreground {
		t.Errorf("Expected foreground %+v, got %+v", top.Foreground, result.Foreground)
	}
}

func TestGlyphBlendingTransparent(t *testing.T) {
	vr := NewViewRenderer(10, 10)

	bottom := Glyph{
		Char:       'A',
		Foreground: Color{R: 100, G: 100, B: 100},
		Alpha:      1.0,
	}

	top := Glyph{
		Char:       'B',
		Foreground: Color{R: 200, G: 200, B: 200},
		Alpha:      0.0, // Fully transparent
		BlendMode:  components.BlendNormal,
	}

	result := vr.blendGlyphs(bottom, top)

	// Should keep bottom unchanged
	if result.Char != 'A' {
		t.Errorf("Expected char 'A', got '%c'", result.Char)
	}

	if result.Foreground != bottom.Foreground {
		t.Errorf("Expected foreground %+v, got %+v", bottom.Foreground, result.Foreground)
	}
}

func TestCharacterSelection(t *testing.T) {
	vr := NewViewRenderer(10, 10)

	// Test character selection at different alpha levels
	bottom := Glyph{Char: 'A'}
	top := Glyph{Char: 'B'}

	// Below threshold - should use bottom
	result := vr.selectChar(bottom, top, 0.3)
	if result != 'A' {
		t.Errorf("Expected 'A' at low alpha, got '%c'", result)
	}

	// Above threshold - should use top
	result = vr.selectChar(bottom, top, 0.7)
	if result != 'B' {
		t.Errorf("Expected 'B' at high alpha, got '%c'", result)
	}

	// Exactly at threshold - should use bottom
	result = vr.selectChar(bottom, top, 0.5)
	if result != 'A' {
		t.Errorf("Expected 'A' at threshold, got '%c'", result)
	}
}

func TestStyleBlending(t *testing.T) {
	vr := NewViewRenderer(10, 10)

	bottom := StyleBold
	top := StyleItalic

	// Below threshold - should keep bottom
	result := vr.blendStyles(bottom, top, 0.3)
	if result != StyleBold {
		t.Errorf("Expected StyleBold at low alpha, got %v", result)
	}

	// Above threshold - should combine
	result = vr.blendStyles(bottom, top, 0.7)
	expected := StyleBold | StyleItalic
	if result != expected {
		t.Errorf("Expected combined styles %v, got %v", expected, result)
	}
}
