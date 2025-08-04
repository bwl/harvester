package systems

import (
	"github.com/charmbracelet/lipgloss"
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"math"
	"math/rand"
	"strconv"
)

type Drawable struct {
	X, Y      int
	Glyph     rune
	Style     lipgloss.Style
	Alpha     float64
	BlendMode components.BlendMode
}

type Render struct{ Output []Drawable }

type Theme struct {
	styles map[components.TileType]lipgloss.Style
}

func getThemeForBiome(biome int) Theme {
	m := map[components.TileType]lipgloss.Style{
		components.TileGalaxy:     lipgloss.NewStyle().Foreground(lipgloss.Color("99")),
		components.TileStar:       lipgloss.NewStyle().Foreground(lipgloss.Color("208")), // orange
		components.TilePlanet:     lipgloss.NewStyle().Foreground(lipgloss.Color("33")),
		components.TileForest:     lipgloss.NewStyle().Foreground(lipgloss.Color("226")), // yellow
		components.TileMountain:   lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
		components.TileRiver:      lipgloss.NewStyle().Foreground(lipgloss.Color("27")), // blue
		components.TileLava:       lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
		components.TileNebula:     lipgloss.NewStyle().Foreground(lipgloss.Color("141")), // magenta
		components.TileGalaxyCore: lipgloss.NewStyle().Foreground(lipgloss.Color("219")).Bold(true),
		components.TileAsteroid:   lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
		components.TileComet:      lipgloss.NewStyle().Foreground(lipgloss.Color("123")),
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

func (r *Render) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	out := r.Output[:0]
	th := getThemeForBiome(ctx.BiomeType)

	// Render tiles with full styling, transparency, and alpha support
	ecs.View2Of[components.Position, components.Tile](w).Each(func(t ecs.Tuple2[components.Position, components.Tile]) {
		style := th.GetStyle(t.B.Type)
		alpha := 1.0 // Default fully opaque for all tiles
		blendMode := components.BlendNormal

		// Check for transparency component (only override if explicitly set)
		if trans, ok := ecs.Get[components.Transparency](w, t.E); ok {
			alpha = trans.Alpha
			blendMode = trans.BlendMode
		}

		out = append(out, Drawable{
			X: int(t.A.X), Y: int(t.A.Y),
			Glyph:     t.B.Glyph,
			Style:     style,
			Alpha:     alpha,
			BlendMode: blendMode,
		})
	})

	// Render entities with full styling, transparency, and alpha support
	ecs.View2Of[components.Position, components.Renderable](w).Each(func(t ecs.Tuple2[components.Position, components.Renderable]) {
		style := th.GetStyle(t.B.TileType)
		alpha := 1.0 // Default fully opaque for all entities
		blendMode := components.BlendNormal

		// Player pulse background using PulseSpring
		if _, isPlayer := ecs.Get[components.Player](w, t.E); isPlayer {
			if ps, ok := ecs.Get[components.PulseSpring](w, t.E); ok {
				c := 255 - int(255*ps.Pos)
				bg := lipgloss.Color(strconv.Itoa(c))
				style = style.Background(bg)
			}
		}

		// Check for transparency component (only override if explicitly set)
		if trans, ok := ecs.Get[components.Transparency](w, t.E); ok {
			alpha = trans.Alpha
			blendMode = trans.BlendMode
		}

		// Apply style modifiers
		if t.B.StyleMod != nil {
			style = applyColorModifier(style, t.B.StyleMod, dt)
		}

		out = append(out, Drawable{
			X: int(t.A.X), Y: int(t.A.Y),
			Glyph:     t.B.Glyph,
			Style:     style,
			Alpha:     alpha,
			BlendMode: blendMode,
		})
	})

	r.Output = out
}
