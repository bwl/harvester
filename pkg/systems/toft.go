package systems

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"math/rand"
)

type WeatherTick struct{}

type Weather struct {
	Rain      bool
	Intensity float64
}

type RiverTile struct{ FlowX, FlowY int }

type TradeRoute struct{ From, To int }

type FactionInfluence struct {
	KingdomID int
	Strength  float64
}

type Wildlife struct{ Hostile bool }

type WildlifeSpawn struct{}

type TradeRoutePatrols struct{}

type Patrol struct{}

type QuestSystem struct{}

type RiverFlow struct{}

type KingdomGuards struct{}

func (s WeatherTick) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	if ctx.CurrentLayer != ecs.LayerPlanetSurface {
		return
	}
	wi, ok := ecs.Get[components.WorldInfo](w, 1)
	if !ok {
		return
	}
	r := rand.New(rand.NewSource(int64(wi.Tick) + 1))
	we, _ := ecs.Get[Weather](w, 1)
	if r.Float64() < 0.1 {
		we.Rain = !we.Rain
		we.Intensity = 0.2 + 0.8*r.Float64()
	}
	ecs.Add(w, 1, we)
}

func (s RiverFlow) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	if ctx.CurrentLayer != ecs.LayerPlanetSurface {
		return
	}
	_ = dt
}

func (s TradeRoutePatrols) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	if ctx.CurrentLayer != ecs.LayerPlanetSurface {
		return
	}
	wi, ok := ecs.Get[components.WorldInfo](w, 1)
	if !ok {
		return
	}
	if wi.Tick%50 == 0 {
		count := 0
		ecs.View2Of[components.Position, Patrol](w).Each(func(t ecs.Tuple2[components.Position, Patrol]) { count++ })
		if count < 5 {
			e := w.Create()
			x, y := wi.Width/2, wi.Height/2
			ecs.Add(w, e, components.Position{X: float64(x), Y: float64(y)})
			ecs.Add(w, e, components.Renderable{Glyph: 'P', TileType: components.TileGalaxy})
			ecs.Add(w, e, Patrol{})
		}
	}
	// wander
	ecs.View2Of[components.Position, Patrol](w).Each(func(t ecs.Tuple2[components.Position, Patrol]) {
		r := rand.Intn(4)
		dx, dy := 0.0, 0.0
		switch r {
		case 0:
			dx = 1
		case 1:
			dx = -1
		case 2:
			dy = 1
		case 3:
			dy = -1
		}
		t.A.X += dx
		t.A.Y += dy
		ecs.Add(w, t.E, *t.A)
	})
}

func (s WildlifeSpawn) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	if ctx.CurrentLayer != ecs.LayerPlanetSurface {
		return
	}
	wi, ok := ecs.Get[components.WorldInfo](w, 1)
	if !ok {
		return
	}
	r := rand.New(rand.NewSource(int64(wi.Tick) + int64(ctx.Depth)*37))
	if r.Float64() < 0.1 {
		e := w.Create()
		x := r.Intn(wi.Width)
		y := r.Intn(wi.Height)
		ecs.Add(w, e, components.Position{X: float64(x), Y: float64(y)})
		ecs.Add(w, e, components.Renderable{Glyph: 'w', TileType: components.TileForest})
		ecs.Add(w, e, Wildlife{Hostile: ctx.Depth > 20})
	}
}

func (q QuestSystem) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	if ctx.CurrentLayer != ecs.LayerPlanetSurface {
		return
	}
	if ctx.QuestProgress.ContractsNeeded == 0 {
		ctx.QuestProgress.ContractsNeeded = 5
	}
	wi, ok := ecs.Get[components.WorldInfo](w, 1)
	if !ok {
		return
	}
	if wi.Tick%100 == 0 && ctx.QuestProgress.ContractsCollected < ctx.QuestProgress.ContractsNeeded {
		ctx.QuestProgress.ContractsCollected++
		if ctx.QuestProgress.ContractsCollected >= ctx.QuestProgress.ContractsNeeded {
			ctx.QuestProgress.RoyalCharterComplete = true
		}
		ecs.SetWorldContext(w, ctx)
	}
}

func (s KingdomGuards) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	if ctx.CurrentLayer != ecs.LayerPlanetSurface {
		return
	}
	_ = dt
}
