package systems

import (
	"harvester/pkg/components"
	"harvester/pkg/data"
	"harvester/pkg/ecs"
)

type Spaceship struct{}

func ApplyDirectionalVelocity(w *ecs.World, e ecs.Entity, dx, dy float64) {
	v, _ := ecs.Get[components.Velocity](w, e)
	v.VX = dx
	v.VY = dy
	ecs.Add(w, e, v)
}

type SpaceObjects struct{}

type FuelSystem struct{}

func (s FuelSystem) Update(dt float64, w *ecs.World) {
	ecs.View2Of[components.FuelTank, components.Velocity](w).Each(func(t ecs.Tuple2[components.FuelTank, components.Velocity]) {
		burn := int((abs(t.B.VX)+abs(t.B.VY))*dt) + 1
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
	ecs.View2Of[components.Position, components.Velocity](w).Each(func(t ecs.Tuple2[components.Position, components.Velocity]) {
		t.A.X += t.B.VX * dt
		t.A.Y += t.B.VY * dt
		ecs.Add(w, t.E, *t.A)
	})
}

type PlanetApproachSystem struct{}

func (s PlanetApproachSystem) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	_ = dt
	// Enter planet when player presses '>' over a planet glyph '1','2','3'
	pressed := false
	ecs.View2Of[EnterPlanet, components.Position](w).Each(func(t ecs.Tuple2[EnterPlanet, components.Position]) { pressed = true })
	if !pressed {
		return
	}
	playerPos := components.Position{}
	ecs.View2Of[components.Player, components.Position](w).Each(func(t ecs.Tuple2[components.Player, components.Position]) {
		playerPos = *t.B
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
