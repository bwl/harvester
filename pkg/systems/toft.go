package systems

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"harvester/pkg/timing"
	"math/rand"
)

type WeatherTick struct{}

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

type RiverFlow struct{}

type KingdomGuards struct{}

func (s WeatherTick) Update(dt float64, w *ecs.World) {
	r := rand.New(rand.NewSource(int64(timing.Tick()) + 1))
	we, _ := ecs.Get[components.Weather](w, 1)
	if r.Float64() < 0.1 {
		we.Rain = !we.Rain
	}
	ecs.Add(w, 1, we)
}

func (s RiverFlow) Update(dt float64, w *ecs.World) {
	_ = dt
}

func (s TradeRoutePatrols) Update(dt float64, w *ecs.World) {
	wi, ok := ecs.Get[components.WorldInfo](w, 1)
	if !ok {
		return
	}
	if timing.Tick()%50 == 0 {
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
	wi, ok := ecs.Get[components.WorldInfo](w, 1)
	if !ok {
		return
	}
	ctx := ecs.GetWorldContext(w)
	r := rand.New(rand.NewSource(int64(timing.Tick()) + int64(ctx.Depth)*37))
	if r.Float64() < 0.1 {
		e := w.Create()
		x := r.Intn(wi.Width)
		y := r.Intn(wi.Height)
		ecs.Add(w, e, components.Position{X: float64(x), Y: float64(y)})
		ecs.Add(w, e, components.Renderable{Glyph: 'w', TileType: components.TileForest})
		ecs.Add(w, e, Wildlife{Hostile: ctx.Depth > 20})
	}
}

func (s KingdomGuards) Update(dt float64, w *ecs.World) {
	_ = dt
}
