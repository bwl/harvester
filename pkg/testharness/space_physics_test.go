package testharness

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"harvester/pkg/engine"
	"harvester/pkg/systems"
	"math"
	"math/rand"
	"testing"
)

func TestRotationSpring(t *testing.T) {
	bs := engine.New(rand.New(rand.NewSource(1)))
	ctx := ecs.GetWorldContext(bs.World)
	ctx.CurrentLayer = ecs.LayerSpace
	ecs.SetWorldContext(bs.World, ctx)
	systems.SetPlayerInput(bs.World, bs.Player, "up")
	systems.SetPlayerInput(bs.World, bs.Player, "right")
	for i := 0; i < 60; i++ {
		bs.Scheduler.Update(1.0/20.0, bs.World)
	}
	spr, _ := ecs.Get[components.SpaceFlightSprings](bs.World, bs.Player)
	if spr.Angle.Pos == 0 {
		t.Fatal("expected angle spring to change with right input")
	}
}

func TestBrakeReducesVelocity(t *testing.T) {
	bs := engine.New(rand.New(rand.NewSource(1)))
	ctx := ecs.GetWorldContext(bs.World)
	ctx.CurrentLayer = ecs.LayerSpace
	ecs.SetWorldContext(bs.World, ctx)
	systems.SetPlayerInput(bs.World, bs.Player, "up")
	for i := 0; i < 60; i++ {
		bs.Scheduler.Update(1.0/20.0, bs.World)
	}
	v0, _ := ecs.Get[components.Velocity](bs.World, bs.Player)
	systems.SetPlayerInput(bs.World, bs.Player, "down")
	for i := 0; i < 30; i++ {
		bs.Scheduler.Update(1.0/20.0, bs.World)
	}
	v1, _ := ecs.Get[components.Velocity](bs.World, bs.Player)
	if math.Abs(v1.VX)+math.Abs(v1.VY) >= math.Abs(v0.VX)+math.Abs(v0.VY) {
		t.Fatal("expected brake to reduce speed")
	}
}
