package systems

import (
	"github.com/charmbracelet/lipgloss"
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"math"
	"math/rand"
	"strconv"
)

type MapRender struct{ Output []Drawable }

type Drawable struct {
	X, Y  int
	Glyph rune
	Style lipgloss.Style
}

type Theme struct {
	styles map[components.TileType]lipgloss.Style
}

func getThemeForBiome(biome int) Theme {
	m := map[components.TileType]lipgloss.Style{
		components.TileGalaxy:   lipgloss.NewStyle().Foreground(lipgloss.Color("99")),
		components.TileStar:     lipgloss.NewStyle().Foreground(lipgloss.Color("208")), // orange
		components.TilePlanet:   lipgloss.NewStyle().Foreground(lipgloss.Color("33")),
		components.TileForest:   lipgloss.NewStyle().Foreground(lipgloss.Color("226")), // yellow
		components.TileMountain: lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
		components.TileRiver:    lipgloss.NewStyle().Foreground(lipgloss.Color("27")), // blue
		components.TileLava:     lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
	}
	return Theme{styles: m}
}

func (t Theme) GetStyle(tt components.TileType) lipgloss.Style {
	if s, ok := t.styles[tt]; ok {
		return s
	}
	return lipgloss.NewStyle()
}

func applyColorModifier(base lipgloss.Style, mod *components.ColorModifier, dt float64) lipgloss.Style {
	switch mod.Special {
	case components.EffectPulsing:
		brightness := 0.5 + 0.5*math.Sin(dt*mod.PulseRate)
		if fg, ok := base.GetForeground().(lipgloss.Color); ok {
			return base.Foreground(adjustBrightness(fg, brightness))
		}
		return base
	case components.EffectTwinkling:
		c := 170 + rand.Intn(81)
		return base.Foreground(lipgloss.Color(strconv.Itoa(c)))
	}
	if mod.TintColor != nil {
		base = base.Foreground(*mod.TintColor)
	}
	return base
}

func adjustBrightness(c lipgloss.Color, factor float64) lipgloss.Color { return c }

func (m *MapRender) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	out := m.Output[:0]
	th := getThemeForBiome(ctx.BiomeType)
	ecs.View2Of[components.Position, components.Tile](w).Each(func(t ecs.Tuple2[components.Position, components.Tile]) {
		style := th.GetStyle(t.B.Type)
		out = append(out, Drawable{X: int(t.A.X), Y: int(t.A.Y), Glyph: t.B.Glyph, Style: style})
	})
	ecs.View2Of[components.Position, components.Renderable](w).Each(func(t ecs.Tuple2[components.Position, components.Renderable]) {
		style := th.GetStyle(t.B.TileType)
		if t.B.StyleMod != nil {
			style = applyColorModifier(style, t.B.StyleMod, dt)
		}
		out = append(out, Drawable{X: int(t.A.X), Y: int(t.A.Y), Glyph: t.B.Glyph, Style: style})
	})
	m.Output = out
}
