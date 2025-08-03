package ui

import (
	"testing"
	"harvester/pkg/rendering"
)

func TestRenderLipgloss_WideRunes(t *testing.T) {
	lines := []string{"こんにちは", "世界🌍"}
	g := rendering.RenderLipglossString(lines, rendering.Color{}, rendering.Color{}, rendering.StyleNone)
	if len(g) != 2 { t.Fatalf("rows=%d", len(g)) }
	if len(g[0]) < 5 { t.Fatalf("expected at least 5 cols, got %d", len(g[0])) }
}
