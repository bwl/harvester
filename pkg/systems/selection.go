package systems

import (
	"fmt"
	"harvester/pkg/components"
	"harvester/pkg/data"
	"harvester/pkg/ecs"
)

type PlanetCard struct {
	Planet data.Planet
	Index  int
}

type PlanetSelection struct{}

func (s PlanetSelection) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	if ctx.CurrentLayer != ecs.LayerSpace {
		return
	}
	// Ensure three planet cards exist once
	created := false
	ecs.View2Of[PlanetCard, PlanetCard](w).Each(func(t ecs.Tuple2[PlanetCard, PlanetCard]) { created = true })
	if !created {
		pg := data.PlanetGenerator{Seed: 1, Biome: data.BiomeToftForest, MaxDepth: 120}
		toft := pg.GenerateToft()
		for i := 0; i < 3; i++ {
			e := w.Create()
			ecs.Add(w, e, PlanetCard{Planet: *toft, Index: i})
			// position cards
			ecs.Add(w, e, components.Position{X: float64(10 + i*20), Y: 5})
			glyph := rune('1' + i)
			ecs.Add(w, e, components.Renderable{Glyph: glyph})
		}
	}
	in, ok := ecs.Get[components.Input](w, 2)
	if ok {
		choice := -1
		if in.Left {
			choice = 0
		}
		if in.Down {
			choice = 1
		}
		if in.Right {
			choice = 2
		}
		if choice >= 0 {
			ecs.SetWorldContext(w, ecs.WorldContext{CurrentLayer: ecs.LayerSpace})
			pg := data.PlanetGenerator{Seed: int64(choice + 1), Biome: data.BiomeToftForest, MaxDepth: 120}
			p := pg.GenerateToft()
			ctx.PlanetID = p.ID
			ctx.Depth = 0
			ctx.BiomeType = int(p.Biome)
			ecs.SetWorldContext(w, ctx)
		}
	}
}

func DescribePlanet(p data.Planet) string {
	return fmt.Sprintf("%s (%d)", p.Name, p.MaxDepth)
}
