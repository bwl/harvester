package ui

import (
	"harvester/pkg/rendering"
	"strings"
)

type uiRightPanel struct {
	g    [][]rendering.Glyph
	w, h int
}

func newUIRightPanel(g [][]rendering.Glyph) *uiRightPanel {
	if g == nil {
		return nil
	}
	return &uiRightPanel{g: g, w: len(g[0]), h: len(g)}
}
func (u *uiRightPanel) GetLayer() rendering.Layer { return rendering.LayerUI }
func (u *uiRightPanel) GetZ() int { return rendering.ZUI }
func (u *uiRightPanel) GetPosition() rendering.Position {
	return rendering.Position{Horizontal: rendering.Right, Vertical: rendering.Top, OffsetY: 2}
}
func (u *uiRightPanel) GetBounds() rendering.Bounds    { return rendering.Bounds{Width: u.w, Height: u.h} }
func (u *uiRightPanel) GetGlyphs() [][]rendering.Glyph { return u.g }

func buildRightPanelGlyphs(m *Model) [][]rendering.Glyph {
	panel := m.renderRightPanel()
	lines := strings.Split(panel, "\n")
	return rendering.RenderLipglossString(lines, rendering.Color{}, rendering.Color{}, rendering.StyleNone)
}
