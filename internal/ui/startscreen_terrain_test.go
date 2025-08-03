package ui

import (
	"harvester/pkg/rendering"
	"strings"
	"testing"
)

func TestStartScreenTerrainBackground(t *testing.T) {
	s := NewStartScreen()
	s.width = 80
	s.height = 24

	// Generate terrain background
	s.generateBackgroundTerrain()
	if s.renderer == nil {
		s.renderer = rendering.NewViewRenderer(s.width, s.height)
	}
	bg := s.renderBackgroundContent()
	if bg == nil {
		t.Fatal("no bg content")
	}
	s.renderer.UnregisterAll()
	s.renderer.RegisterContent(bg)
	background := s.renderer.Render()

	if background == "" {
		t.Error("Background terrain should not be empty")
	}

	// Check that background contains terrain elements
	if !strings.Contains(background, "#") && !strings.Contains(background, "~") {
		t.Error("Background should contain forest (#) or river (~) tiles")
	}

	// Check dimensions
	lines := strings.Split(background, "\n")
	if len(lines) < s.height {
		t.Errorf("Background should have at least %d lines, got %d", s.height, len(lines))
	}
}

// Deprecated: overlay logic removed; compositor is used instead
func TestStartScreenOverlay(t *testing.T) {
	 t.Skip("overlay removed")
	s := NewStartScreen()
	s.width = 80
	s.height = 24

	// Generate terrain background
	s.generateBackgroundTerrain()
	if s.renderer == nil {
		s.renderer = rendering.NewViewRenderer(s.width, s.height)
	}
	bg := s.renderBackgroundContent()
	if bg == nil {
		t.Fatal("no bg content")
	}
	s.renderer.UnregisterAll()
	s.renderer.RegisterContent(bg)
	_ = s.renderer.Render()

	// Get menu content
	_ = s.renderMainMenu()

	// overlay removed; nothing to assert here
	return
}

func TestANSIStripper(t *testing.T) {
	input := "\x1b[31mHello\x1b[0m World"
	expected := "Hello World"
	result := stripANSI(input)

	if result != expected {
		t.Errorf("stripANSI failed: expected '%s', got '%s'", expected, result)
	}
}
