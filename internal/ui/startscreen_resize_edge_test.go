package ui

import (
	"testing"
	tea "github.com/charmbracelet/bubbletea"
)

func TestStartScreen_Resize_87x342(t *testing.T) {
	s := NewStartScreen()
	_, _ = s.Update(tea.WindowSizeMsg{Width: 342, Height: 87})
	out := s.View()
	if out == "" || out == "Loading..." { t.Fatal("no output") }
}
