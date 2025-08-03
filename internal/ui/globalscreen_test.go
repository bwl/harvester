package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/rendering"
	"harvester/pkg/timing"
)

func TestGlobalScreenCreation(t *testing.T) {
	gs := NewGlobalScreen()
	if gs == nil {
		t.Fatal("NewGlobalScreen should return a valid global screen")
	}

	if gs.currentScreen != ScreenStart {
		t.Error("GlobalScreen should start with StartScreen")
	}

	if gs.subScreen == nil {
		t.Error("GlobalScreen should have a valid sub-screen")
	}
}

func TestGlobalScreenShutdownAnimation(t *testing.T) {
	gs := NewGlobalScreen()
	gs.width = 80
	gs.height = 24

	// Start shutdown animation
	gs.startShutdownAnimation()

	if gs.shutdownAnim == nil {
		t.Error("Shutdown animation should be initialized")
	}

	if !gs.quitting {
		t.Error("Quitting flag should be set")
	}
}

func TestGlobalScreenRenderableInterface(t *testing.T) {
	gs := NewGlobalScreen()
	gs.width = 80
	gs.height = 10

	// Test that GlobalScreen implements RenderableScreen
	var renderable RenderableScreen = gs
	if renderable == nil {
		t.Error("GlobalScreen should implement RenderableScreen interface")
	}

	// Test with actual ViewRenderer  
	renderer := rendering.NewViewRenderer(80, 10)
	
	// Should not panic when registering content
	gs.RegisterContent(renderer)

	// Test with shutdown animation
	gs.shutdownAnim = timing.NewAnimation(10, false)
	gs.RegisterContent(renderer)
	
	// Should complete without errors (overlay registration tested elsewhere)
}

func TestGlobalScreenWindowResize(t *testing.T) {
	gs := NewGlobalScreen()

	// Test window resize message
	windowMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
	newModel, _ := gs.Update(windowMsg)

	updatedGS := newModel.(*GlobalScreen)
	if updatedGS.width != 120 || updatedGS.height != 40 {
		t.Error("GlobalScreen should update dimensions on window resize")
	}
}

func TestGlobalScreenQuitKeys(t *testing.T) {
	gs := NewGlobalScreen()
	gs.width = 80
	gs.height = 24

	// Test Ctrl+C (immediate quit)
	ctrlCMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := gs.Update(ctrlCMsg)

	if cmd == nil {
		t.Error("Ctrl+C should return a quit command")
	}

	// Test Q key (animated quit, but only when not in game)
	qMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	gs.currentScreen = ScreenStart // Ensure we're not in game
	newModel, _ := gs.Update(qMsg)

	updatedGS := newModel.(*GlobalScreen)
	if updatedGS.shutdownAnim == nil {
		t.Error("Q key should start shutdown animation when not in game")
	}
}

func TestGlobalScreenGameTransition(t *testing.T) {
	gs := NewGlobalScreen()

	// Simulate start screen result
	result := &StartResult{Action: ActionNewGame}
	newModel, _ := gs.handleStartScreenResult(result)

	updatedGS := newModel.(*GlobalScreen)
	if !updatedGS.transitioning {
		t.Error("Should be transitioning after start screen result")
	}

	if updatedGS.nextScreen != ScreenSpace {
		t.Error("Should be transitioning to space screen")
	}
}

func TestGlobalScreenIgnoreInputDuringAnimation(t *testing.T) {
	gs := NewGlobalScreen()
	gs.width = 80
	gs.height = 24

	// Start animation
	gs.startShutdownAnimation()

	// Try to send input
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
	newModel, _ := gs.Update(keyMsg)

	// Should return same model without processing input
	if newModel != gs {
		t.Error("Input should be ignored during shutdown animation")
	}
}

func TestLineSplitAndJoin(t *testing.T) {
	content := "Line 1\nLine 2\nLine 3"
	lines := splitLines(content)

	expectedLines := []string{"Line 1", "Line 2", "Line 3"}
	if len(lines) != len(expectedLines) {
		t.Errorf("Expected %d lines, got %d", len(expectedLines), len(lines))
	}

	for i, line := range lines {
		if line != expectedLines[i] {
			t.Errorf("Line %d: expected %s, got %s", i, expectedLines[i], line)
		}
	}

	// Test join
	joined := joinLines(lines)
	if joined != content {
		t.Errorf("Join should restore original content: expected %s, got %s", content, joined)
	}
}

func TestEmptyContentHandling(t *testing.T) {
	// Test empty content
	lines := splitLines("")
	if len(lines) != 0 {
		t.Error("Empty content should produce empty lines slice")
	}

	joined := joinLines([]string{})
	if joined != "" {
		t.Error("Empty lines should produce empty string")
	}
}
