package systems

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

type SurfaceHeartbeat struct{}

type DepthProgression struct{}

type TerrainGen struct{}

type RiverTag struct{}

type SurfaceMovement struct{}

func (s SurfaceHeartbeat) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	if ctx.CurrentLayer != ecs.LayerPlanetSurface {
		return
	}
	wi, ok := ecs.Get[components.WorldInfo](w, 1)
	if !ok {
		return
	}
	wi.Tick++
	ecs.Add(w, 1, wi)
}

func (t TerrainGen) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	if ctx.CurrentLayer != ecs.LayerPlanetSurface {
		return
	}
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
		ecs.Add(w, e, RiverTag{})
	}
}

func (d DepthProgression) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	if ctx.CurrentLayer != ecs.LayerPlanetSurface {
		return
	}
	in, ok := ecs.Get[components.Input](w, 2)
	if !ok {
		return
	}
	we, _ := ecs.Get[Weather](w, 1)
	step := 1
	if we.Rain {
		step = 2
	}
	playerPos, _ := ecs.Get[components.Position](w, 2)
	river := false
	ecs.View2Of[components.Position, RiverTag](w).Each(func(t ecs.Tuple2[components.Position, RiverTag]) {
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
	in, ok := ecs.Get[components.Input](w, 2)
	if !ok {
		return
	}
	p, ok := ecs.Get[components.Position](w, 2)
	if !ok {
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
	we, _ := ecs.Get[Weather](w, 1)
	speed := 1.0
	if we.Rain {
		speed *= 0.5
	}
	onRiver := false
	ecs.View2Of[components.Position, RiverTag](w).Each(func(t ecs.Tuple2[components.Position, RiverTag]) {
		if int(t.A.X) == int(p.X) && int(t.A.Y) == int(p.Y) {
			onRiver = true
		}
	})
	if onRiver {
		speed *= 0.5
	}
	p.X += dx * speed
	p.Y += dy * speed
	ecs.Add(w, 2, p)
	_ = ctx
}
