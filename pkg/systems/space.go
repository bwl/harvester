package systems

import (
	"github.com/charmbracelet/harmonica"
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
	const angW, angZ = 6.0, 0.6
	const thrW, thrZ = 5.0, 0.7
	const velW, velZ = 6.0, 0.6
	spAngle := harmonica.NewSpring(harmonica.FPS(60), angW, angZ)
	spThrust := harmonica.NewSpring(harmonica.FPS(60), thrW, thrZ)
	spVel := harmonica.NewSpring(harmonica.FPS(60), velW, velZ)
	ecs.View2Of[components.Position, components.Velocity](w).Each(func(t ecs.Tuple2[components.Position, components.Velocity]) {
		// fetch springs and input/orientation
		spr, _ := ecs.Get[components.SpaceFlightSprings](w, t.E)
		inp, _ := ecs.Get[components.Input](w, t.E)
		// update springs
		spr.Angle.Pos, spr.Angle.Vel = spAngle.Update(spr.Angle.Pos, spr.Angle.Vel, spr.Angle.Target)
		spr.Thrust.Pos, spr.Thrust.Vel = spThrust.Update(spr.Thrust.Pos, spr.Thrust.Vel, spr.Thrust.Target)
		if inp.Down {
			spr.VelX.Target, spr.VelY.Target = 0, 0
		}
		spr.VelX.Pos, spr.VelX.Vel = spVel.Update(spr.VelX.Pos, spr.VelX.Vel, spr.VelX.Target)
		spr.VelY.Pos, spr.VelY.Vel = spVel.Update(spr.VelY.Pos, spr.VelY.Vel, spr.VelY.Target)
		// physics integration: accel from thrust+angle
		ax := spr.Thrust.Pos * cos(spr.Angle.Pos)
		ay := spr.Thrust.Pos * sin(spr.Angle.Pos)
		vx := t.B.VX + ax*dt
		vy := t.B.VY + ay*dt
		// apply velocity springs as damping toward targets
		vx, _ = spVel.Update(vx, 0, spr.VelX.Pos)
		vy, _ = spVel.Update(vy, 0, spr.VelY.Pos)
		t.A.X += vx * dt
		t.A.Y += vy * dt
		t.B.VX, t.B.VY = vx, vy
		ecs.Add(w, t.E, *t.A)
		ecs.Add(w, t.E, *t.B)
		ecs.Add(w, t.E, spr)
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
