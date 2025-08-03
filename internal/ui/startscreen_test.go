package ui

import (
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestStartScreenCreation(t *testing.T) {
	s := NewStartScreen()
	if s == nil {
		t.Fatal("NewStartScreen should return a valid start screen")
	}

	// Should have at least "New Game" and "Quit" options
	if len(s.menuItems) < 2 {
		t.Error("Start screen should have at least 2 menu items")
	}

	// Should contain "New Game"
	found := false
	for _, item := range s.menuItems {
		if item == "New Game" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Start screen should always include 'New Game' option")
	}
}

func TestStartScreenNavigation(t *testing.T) {
	s := NewStartScreen()
	s.width = 80
	s.height = 24

	if len(s.menuItems) < 2 {
		t.Skip("Need at least 2 menu items to test navigation")
	}

	initialSelected := s.selected

	// Test navigation down
	newModel, _ := s.updateMainMenu(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	newS := newModel.(*StartScreen)
	if newS.selected <= initialSelected {
		t.Error("Down navigation should increase selection")
	}

	// Test navigation up
	s.selected = 1 // Set to second item
	newModel, _ = s.updateMainMenu(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	newS = newModel.(*StartScreen)
	if newS.selected != 0 {
		t.Error("Up navigation should decrease selection")
	}
}

func TestSaveFileDetection(t *testing.T) {
	s := NewStartScreen()

	// The saveSlots should be initialized (even if files don't exist)
	if len(s.saveSlots) != 3 {
		t.Error("Should detect 3 save slots")
	}

	for i, slot := range s.saveSlots {
		if slot.SlotNum != i+1 {
			t.Errorf("Slot %d should have SlotNum %d", i, i+1)
		}
	}
}

func TestMenuItemsWithSaves(t *testing.T) {
	// Create a temporary autosave for testing
	_ = os.MkdirAll(".saves", 0o755)
	_ = os.WriteFile(".saves/autosave.gz", []byte("test"), 0o644)
	defer os.Remove(".saves/autosave.gz")

	s := NewStartScreen()

	// Should have "Continue" option when autosave exists
	found := false
	for _, item := range s.menuItems {
		if item == "Continue" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Should include 'Continue' option when autosave exists")
	}
}

func TestQuitAction(t *testing.T) {
	s := NewStartScreen()
	s.width = 80
	s.height = 24

	// Test quit with 'q'
	newModel, cmd := s.updateMainMenu(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	startScreen := newModel.(*StartScreen)

	if startScreen.result == nil || startScreen.result.Action != ActionQuit {
		t.Error("Pressing 'q' should set quit action")
	}

	if cmd == nil {
		t.Error("Quit should return a quit command")
	}
}

func TestRenderWithoutCrash(t *testing.T) {
	s := NewStartScreen()
	s.width = 80
	s.height = 24

	// Should not crash when rendering
	view := s.View()
	if view == "" {
		t.Error("View should return content")
	}

	// Test slot view rendering
	s.showSlots = true
	slotView := s.View()
	if slotView == "" {
		t.Error("Slot view should return content")
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{500, "500B"},
		{1500, "1.5KB"},
		{1500000, "1.4MB"},
	}

	for _, test := range tests {
		result := formatFileSize(test.size)
		if result != test.expected {
			t.Errorf("formatFileSize(%d) = %s, expected %s", test.size, result, test.expected)
		}
	}
}
