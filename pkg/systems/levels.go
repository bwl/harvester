package systems

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

type LevelTag struct{ ID int }

type LevelManager struct{}

func (LevelManager) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	if ctx.CurrentLayer == ecs.LayerPlanetSurface {
		// destroy space visuals: stars, planet cards
		toDestroy := make([]ecs.Entity, 0, 128)
		ecs.View2Of[components.Tile, components.Position](w).Each(func(t ecs.Tuple2[components.Tile, components.Position]) {
			if t.A.Glyph == '*' || (t.A.Glyph >= '1' && t.A.Glyph <= '3') {
				toDestroy = append(toDestroy, t.E)
			}
		})
		ecs.View2Of[components.Renderable, components.Position](w).Each(func(t ecs.Tuple2[components.Renderable, components.Position]) {
			if t.A.Glyph == '*' || (t.A.Glyph >= '1' && t.A.Glyph <= '3') {
				toDestroy = append(toDestroy, t.E)
			}
		})
		for _, e := range toDestroy {
			w.Destroy(e)
		}
	}
}
