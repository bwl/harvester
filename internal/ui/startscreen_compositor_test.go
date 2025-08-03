package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/rendering"
)

func TestStartScreen_CompositorRendersMenuAndBackground(t *testing.T) {
	s := NewStartScreen()
	// simulate size
	_, _ = s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	r := rendering.NewViewRenderer(80, 24)
	s.RegisterContent(r)
	out := r.Render()
	if len(out) == 0 {
		t.Fatal("expected rendered output")
	}
	if !strings.Contains(out, "BUBBLE ROUGE") {
		t.Error("menu title not present in output")
	}
	if len(strings.Split(out, "\n")) < 10 {
		t.Error("output too short, expected screenful of content")
	}
}

func TestLipglossRaster_Basic(t *testing.T) {
	g := rendering.RenderLipglossString([]string{"Hello"}, rendering.Color{}, rendering.Color{}, rendering.StyleNone)
	if len(g) != 1 || len(g[0]) != 5 {
		t.Fatalf("unexpected glyph dims %dx%d", len(g), len(g[0]))
	}
}
