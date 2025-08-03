package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/rendering"
)

func TestStartScreen_ResizeTable(t *testing.T) {
	s := NewStartScreen()
	cases := []struct{ w, h int }{
		{80, 24}, {100, 30}, {60, 20}, {120, 40},
	}
	for _, c := range cases {
		_, _ = s.Update(tea.WindowSizeMsg{Width: c.w, Height: c.h})
		r := rendering.NewViewRenderer(c.w, c.h)
		s.RegisterContent(r)
		out := r.Render()
		if out == "" {
			t.Fatalf("no output for %dx%d", c.w, c.h)
		}
	}
}
