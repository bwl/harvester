package ui

import (
	"harvester/pkg/components"
	"harvester/pkg/rendering"
	"testing"
)

func TestFadeOverlay(t *testing.T) {
	// Test fade in
	fadeIn := NewFadeOverlay(10, 5, 0.5, true)
	glyphs := fadeIn.GetGlyphs()

	if len(glyphs) != 5 {
		t.Errorf("Expected 5 rows, got %d", len(glyphs))
	}

	if len(glyphs[0]) != 10 {
		t.Errorf("Expected 10 columns, got %d", len(glyphs[0]))
	}

	// Test that fade in works differently than fade out
	fadeInTest := NewFadeOverlay(10, 5, 0.5, true)
	fadeOutTest := NewFadeOverlay(10, 5, 0.5, false)

	glyphsIn := fadeInTest.GetGlyphs()
	glyphsOut := fadeOutTest.GetGlyphs()

	// Fade in and fade out should produce different alpha values at same progress
	if glyphsIn[0][0].Alpha == glyphsOut[0][0].Alpha {
		t.Error("Fade in and fade out should produce different alpha values")
	}
}

func TestPulseOverlay(t *testing.T) {
	pulse := NewPulseOverlay(5, 3, 0.0, 0.5)
	glyphs := pulse.GetGlyphs()

	if len(glyphs) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(glyphs))
	}

	// Test that pulse overlay uses additive blending
	if glyphs[0][0].BlendMode != components.BlendAdditive {
		t.Error("Pulse overlay should use additive blending")
	}

	// Test different time values produce different alpha
	pulse1 := NewPulseOverlay(5, 3, 0.0, 0.5)
	pulse2 := NewPulseOverlay(5, 3, 0.25, 0.5) // Quarter cycle later

	glyphs1 := pulse1.GetGlyphs()
	glyphs2 := pulse2.GetGlyphs()

	// Should have different alpha values due to sine wave
	if glyphs1[0][0].Alpha == glyphs2[0][0].Alpha {
		t.Error("Pulse overlay should vary alpha over time")
	}
}

func TestWaveDistortionOverlay(t *testing.T) {
	wave := NewWaveDistortionOverlay(8, 6, 0.0, 0.3, 2.0)
	glyphs := wave.GetGlyphs()

	if len(glyphs) != 6 {
		t.Errorf("Expected 6 rows, got %d", len(glyphs))
	}

	// Test that wave uses multiply blending
	if glyphs[0][0].BlendMode != components.BlendMultiply {
		t.Error("Wave distortion should use multiply blending")
	}

	// Test that different rows have different alpha values (wave effect)
	alphas := make([]float64, len(glyphs))
	for i, row := range glyphs {
		alphas[i] = row[0].Alpha
	}

	// Check that not all alphas are the same (wave creates variation)
	allSame := true
	for i := 1; i < len(alphas); i++ {
		if alphas[i] != alphas[0] {
			allSame = false
			break
		}
	}

	if allSame {
		t.Error("Wave distortion should create varying alpha across rows")
	}
}

func TestSpringBounceOverlay(t *testing.T) {
	bounce := NewSpringBounceOverlay(6, 10, 0.5, 8)
	glyphs := bounce.GetGlyphs()

	if len(glyphs) != 10 {
		t.Errorf("Expected 10 rows, got %d", len(glyphs))
	}

	// Test that bounce effect creates a line-like pattern
	// Count non-zero alpha positions
	nonZeroAlphaCount := 0
	for _, row := range glyphs {
		if row[0].Alpha > 0.0 {
			nonZeroAlphaCount++
		}
	}

	// Should have a focused bounce effect (not covering entire height)
	if nonZeroAlphaCount >= len(glyphs) {
		t.Error("Bounce effect should be localized, not cover entire height")
	}

	// Test that bounce uses block character
	for _, row := range glyphs {
		if row[0].Alpha > 0.0 && row[0].Char != '▬' {
			t.Error("Bounce effect should use block character '▬'")
		}
	}
}

func TestOverlayInterfaces(t *testing.T) {
	// Test that all overlays implement RenderableContent interface
	var overlays []rendering.RenderableContent

	// This will fail to compile if interfaces aren't properly implemented
	fade := NewFadeOverlay(5, 5, 0.5, true)
	pulse := NewPulseOverlay(5, 5, 0.0, 0.5)
	wave := NewWaveDistortionOverlay(5, 5, 0.0, 0.3, 1.0)
	bounce := NewSpringBounceOverlay(5, 5, 0.5, 3)

	overlays = append(overlays, fade, pulse, wave, bounce)

	if len(overlays) != 4 {
		t.Error("All overlays should implement RenderableContent interface")
	}
}

func TestHarmonicaEffectConsistency(t *testing.T) {
	// Test that harmonica effects are deterministic for same inputs
	fade1 := NewFadeOverlay(5, 5, 0.7, false)
	fade2 := NewFadeOverlay(5, 5, 0.7, false)

	glyphs1 := fade1.GetGlyphs()
	glyphs2 := fade2.GetGlyphs()

	// Should produce identical results for identical inputs
	if glyphs1[0][0].Alpha != glyphs2[0][0].Alpha {
		t.Error("Harmonica effects should be deterministic for same inputs")
	}
}
