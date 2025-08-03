package systems

import (
	"harvester/pkg/components"
	"harvester/pkg/data"
	"harvester/pkg/ecs"
)

type Spaceship struct{}

type FuelTank struct{ Current int }

type Velocity struct{ X, Y float64 }

func ApplyDirectionalVelocity(w *ecs.World, e ecs.Entity, dx, dy float64) {
	v, _ := ecs.Get[Velocity](w, e)
	v.X = dx
	v.Y = dy
	ecs.Add(w, e, v)
}

type Planet struct {
	ID   int
	Name string
}

type SpaceObjects struct{}

type FuelSystem struct{}

func (s FuelSystem) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	if ctx.CurrentLayer != ecs.LayerSpace {
		return
	}
	ecs.View2Of[FuelTank, Velocity](w).Each(func(t ecs.Tuple2[FuelTank, Velocity]) {
		burn := int((abs(t.B.X)+abs(t.B.Y))*dt) + 1
		if t.A.Current > 0 {
			t.A.Current -= burn
			if t.A.Current < 0 {
				t.A.Current = 0
			}
			ecs.Add(w, t.E, *t.A)
		}
	})
}

type SpaceMovement struct{}

func (s SpaceMovement) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	if ctx.CurrentLayer != ecs.LayerSpace {
		return
	}
	ecs.View2Of[components.Position, Velocity](w).Each(func(t ecs.Tuple2[components.Position, Velocity]) {
		t.A.X += t.B.X * dt
		t.A.Y += t.B.Y * dt
		ecs.Add(w, t.E, *t.A)
	})
}

type PlanetApproachSystem struct{}

func (s PlanetApproachSystem) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	if ctx.CurrentLayer != ecs.LayerSpace {
		return
	}
	_ = dt
	// Enter planet when player presses '>' over a planet glyph '1','2','3'
	pressed := false
	ecs.View2Of[EnterPlanet, components.Position](w).Each(func(t ecs.Tuple2[EnterPlanet, components.Position]) { pressed = true })
	if !pressed {
		return
	}
	playerPos := components.Position{}
	ecs.View2Of[components.Position, components.Renderable](w).Each(func(t ecs.Tuple2[components.Position, components.Renderable]) {
		if t.B.Glyph == '@' {
			playerPos = *t.A
		}
	})
	enterID := -1
	ecs.View2Of[components.Position, components.Renderable](w).Each(func(t ecs.Tuple2[components.Position, components.Renderable]) {
		if t.B.Glyph >= '1' && t.B.Glyph <= '3' {
			if int(t.A.X) == int(playerPos.X) && int(t.A.Y) == int(playerPos.Y) {
				enterID = int(t.B.Glyph - '0')
			}
		}
	})
	if enterID > 0 {
		pg := data.PlanetGenerator{Seed: int64(enterID), Biome: data.BiomeToftForest, MaxDepth: 120}
		p := pg.GenerateToft()
		ctx.CurrentLayer = ecs.LayerPlanetSurface
		ctx.PlanetID = p.ID
		ctx.Depth = 0
		ecs.SetWorldContext(w, ctx)
	}
}

func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
