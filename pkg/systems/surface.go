package systems

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

type SurfaceHeartbeat struct{}

type DepthProgression struct{}

type TerrainGen struct{}

type SurfaceMovement struct{}

func (s SurfaceHeartbeat) Update(dt float64, w *ecs.World) {
	wi, ok := ecs.Get[components.WorldInfo](w, 1)
	if !ok {
		return
	}
	// Note: Tick is now handled by global timing system
	ecs.Add(w, 1, wi)
}

func (t TerrainGen) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	wi, ok := ecs.Get[components.WorldInfo](w, 1)
	if !ok {
		return
	}
	seed := int64(ctx.PlanetID*100000 + ctx.Depth)
	r := ecs.RandFromSeed(seed)
	for y := 0; y < wi.Height; y++ {
		for x := 0; x < wi.Width; x++ {
			if (x+y+int(r.Int63()%7))%17 == 0 {
				e := w.Create()
				ecs.Add(w, e, components.Position{X: float64(x), Y: float64(y)})
				ecs.Add(w, e, components.Tile{Glyph: '#', Type: components.TileForest})
			}
		}
	}
	// simple river: vertical line
	x0 := wi.Width/3 + int(r.Int63()%5)
	for y := 0; y < wi.Height; y++ {
		e := w.Create()
		ecs.Add(w, e, components.Position{X: float64(x0), Y: float64(y)})
		ecs.Add(w, e, components.Tile{Glyph: '~', Type: components.TileRiver})
		ecs.Add(w, e, components.RiverTag{})

		// Make river slightly transparent to show terrain underneath
		ecs.Add(w, e, components.Transparency{
			Alpha:     0.8, // 80% opacity - slightly see-through
			BlendMode: components.BlendNormal,
		})
	}

	// Add some fog patches for atmosphere
	fogCount := int(r.Int63() % 5) // 0-4 fog patches
	for i := 0; i < fogCount; i++ {
		fx := int(r.Int63() % int64(wi.Width))
		fy := int(r.Int63() % int64(wi.Height))

		// Create a small fog patch (3x3 area)
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				x, y := fx+dx, fy+dy
				if x >= 0 && x < wi.Width && y >= 0 && y < wi.Height {
					// Random chance for each fog cell
					if r.Float64() < 0.6 {
						e := w.Create()
						ecs.Add(w, e, components.Position{X: float64(x), Y: float64(y)})
						ecs.Add(w, e, components.Tile{Glyph: '░', Type: components.TileForest})
						ecs.Add(w, e, components.Transparency{
							Alpha:     0.2 + r.Float64()*0.3, // 20-50% opacity
							BlendMode: components.BlendNormal,
						})
					}
				}
			}
		}
	}
}

func (d DepthProgression) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	var in *components.Input
	var playerPos components.Position
	ecs.View3Of[components.Player, components.Input, components.Position](w).Each(func(t ecs.Tuple3[components.Player, components.Input, components.Position]) {
		in = t.B
		playerPos = *t.C
	})
	if in == nil {
		return
	}
	we, _ := ecs.Get[components.Weather](w, 1)
	step := 1
	if we.Rain {
		step = 2
	}
	river := false
	ecs.View2Of[components.Position, components.RiverTag](w).Each(func(t ecs.Tuple2[components.Position, components.RiverTag]) {
		if int(t.A.X) == int(playerPos.X) && int(t.A.Y) == int(playerPos.Y) {
			river = true
		}
	})
	if river {
		step++
	}
	if in.Down && ctx.Depth < 10000 {
		ctx.Depth += step
	}
	if in.Up && ctx.Depth > 0 {
		ctx.Depth -= step
		if ctx.Depth < 0 {
			ctx.Depth = 0
		}
	}
	ecs.SetWorldContext(w, ctx)
}

func (s SurfaceMovement) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	if ctx.CurrentLayer != ecs.LayerPlanetSurface {
		return
	}
	var in *components.Input
	var p *components.Position
	ecs.View3Of[components.Player, components.Input, components.Position](w).Each(func(t ecs.Tuple3[components.Player, components.Input, components.Position]) {
		in = t.B
		p = t.C
	})
	if in == nil || p == nil {
		return
	}
	dx, dy := 0.0, 0.0
	if in.Left {
		dx = -1
	}
	if in.Right {
		dx = 1
	}
	if in.Up {
		dy = -1
	}
	if in.Down {
		dy = 1
	}
	if dx == 0 && dy == 0 {
		return
	}
	we, _ := ecs.Get[components.Weather](w, 1)
	speed := 1.0
	if we.Rain {
		speed *= 0.5
	}
	onRiver := false
	ecs.View2Of[components.Position, components.RiverTag](w).Each(func(t ecs.Tuple2[components.Position, components.RiverTag]) {
		if int(t.A.X) == int(p.X) && int(t.A.Y) == int(p.Y) {
			onRiver = true
		}
	})
	if onRiver {
		speed *= 0.5
	}
	p.X += dx * speed
	p.Y += dy * speed
	ecs.View1Of[components.Player](w).Each(func(e ecs.Entity, _ *components.Player) { ecs.Add(w, e, *p) })
	_ = ctx
}
