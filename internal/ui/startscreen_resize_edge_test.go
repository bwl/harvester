package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/rendering"
	"testing"
)

func TestStartScreen_Resize_87x342(t *testing.T) {
	s := NewStartScreen()
	_, _ = s.Update(tea.WindowSizeMsg{Width: 342, Height: 87})
	r := rendering.NewViewRenderer(342, 87)
	s.RegisterContent(r)
	out := r.Render()
	if out == "" {
		t.Fatal("no output")
	}
}
