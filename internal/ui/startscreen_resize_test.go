package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestStartScreen_ResizeTable(t *testing.T) {
	s := NewStartScreen()
	cases := []struct{ w,h int }{
		{80,24}, {100,30}, {60,20}, {120,40},
	}
	for _, c := range cases {
		_, _ = s.Update(tea.WindowSizeMsg{Width: c.w, Height: c.h})
		out := s.View()
		if out == "" || out == "Loading..." { t.Fatalf("no output for %dx%d", c.w, c.h) }
	}
}
