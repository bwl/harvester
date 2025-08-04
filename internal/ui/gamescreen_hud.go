package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/v2"
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"harvester/pkg/rendering"
)

func buildGameGlyphs(m *Model, w, h int) [][]rendering.Glyph {
	if w <= 0 || h <= 0 {
		return nil
	}
	// Use unified render system for all drawables
	m.render.Update(0, m.world)
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
	for _, d := range m.render.Output {
		x := d.X - mx0
		y := d.Y - my0
		if x >= 0 && y >= 0 && x < w && y < h {
			glyphs[y][x] = rendering.Glyph{Char: rune(d.Glyph)}
		}
	}
	return glyphs
}

type hudContent struct {
	model *Model
	w     int
}

func newHUDContent(model *Model) *hudContent {
	return &hudContent{model: model, w: 80} // Default width
}

func (h *hudContent) GetLayer() rendering.Layer { return rendering.LayerMenu }
func (h *hudContent) GetZ() int                 { return rendering.ZHUD }

func (h *hudContent) ToLipglossLayer() *lipgloss.Layer {
	if h.model == nil {
		return lipgloss.NewLayer("").X(0).Y(2).Z(h.GetZ()).ID("hud")
	}

	ps, _ := ecs.Get[components.PlayerStats](h.model.world, h.model.player)
	ctx := ecs.GetWorldContext(h.model.world)
	hudText := fmt.Sprintf("HP:%d Fuel:%d Drive:%d  Layer:%s  Tick:%d",
		ps.Hull, ps.Fuel, ps.Drive, layerName(ctx.CurrentLayer), int(h.model.frame))

	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff"))
	styledText := style.Render(hudText)

	return lipgloss.NewLayer(styledText).
		X(0).
		Y(2).
		Z(h.GetZ()).
		ID("hud")
}
