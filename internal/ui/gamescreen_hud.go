package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"harvester/pkg/rendering"
	"harvester/pkg/systems"
)

func buildGameGlyphs(m *Model, w, h int) [][]rendering.Glyph {
	if w <= 0 || h <= 0 {
		return nil
	}
	m.render.Update(0, m.world)
	mr := systems.MapRender{}
	mr.Update(0, m.world)
	cam, _ := ecs.Get[components.Camera](m.world, m.player)
	mx0, my0 := cam.X, cam.Y
	glyphs := make([][]rendering.Glyph, h)
	for y := 0; y < h; y++ {
		row := make([]rendering.Glyph, w)
		for x := 0; x < w; x++ {
			row[x] = rendering.Glyph{Char: '.'}
		}
		glyphs[y] = row
	}
	for _, d := range mr.Output {
		x := d.X - mx0
		y := d.Y - my0
		if x >= 0 && y >= 0 && x < w && y < h {
			glyphs[y][x] = rendering.Glyph{Char: rune(d.Glyph)}
		}
	}
	for _, d := range m.render.Output {
		x := d.X - mx0
		y := d.Y - my0
		if x >= 0 && y >= 0 && x < w && y < h {
			glyphs[y][x] = rendering.Glyph{Char: rune(d.Glyph)}
		}
	}
	return glyphs
}

func buildHUDGlyphs(m *Model, w int) [][]rendering.Glyph {
	ps, _ := ecs.Get[components.PlayerStats](m.world, m.player)
	ctx := ecs.GetWorldContext(m.world)
	hudText := fmt.Sprintf("HP:%d Fuel:%d Drive:%d  Layer:%s  Tick:%d",
		ps.Hull, ps.Fuel, ps.Drive, layerName(ctx.CurrentLayer), int(m.frame))
	style := lipgloss.NewStyle().Bold(true)
	line := style.Render(hudText)
	return rendering.RenderLipglossString([]string{strings.TrimRight(line, "\n")}, rendering.Color{}, rendering.Color{}, rendering.StyleNone)
}


type hudContent struct {
	g [][]rendering.Glyph
	w int
}

func newHUDContent(g [][]rendering.Glyph) *hudContent {
	if g == nil {
		return nil
	}
	return &hudContent{g: g, w: len(g[0])}
}
func (h *hudContent) GetLayer() rendering.Layer { return rendering.LayerMenu }
func (h *hudContent) GetZ() int                 { return rendering.ZHUD }
func (h *hudContent) GetPosition() rendering.Position {
	return rendering.Position{Horizontal: rendering.Left, Vertical: rendering.Top, OffsetY: 2}
}
func (h *hudContent) GetBounds() rendering.Bounds    { return rendering.Bounds{Width: h.w, Height: 1} }
func (h *hudContent) GetGlyphs() [][]rendering.Glyph { return h.g }
