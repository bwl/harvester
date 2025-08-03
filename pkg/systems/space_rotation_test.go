package systems

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"math/rand"
	"testing"
)

func TestSpaceRotationSpring(t *testing.T) {
	bs := newTestBootstrap(rand.New(rand.NewSource(1)))
	ctx := ecs.GetWorldContext(bs.World)
	ctx.CurrentLayer = ecs.LayerSpace
	ecs.SetWorldContext(bs.World, ctx)
	SetPlayerInput(bs.World, bs.Player, "up")
	SetPlayerInput(bs.World, bs.Player, "right")
	for i := 0; i < 60; i++ {
		bs.Scheduler.Update(1.0/20.0, bs.World)
	}
	spr, _ := ecs.Get[components.SpaceFlightSprings](bs.World, bs.Player)
	if spr.Angle.Pos == 0 {
		t.Fatal("expected angle spring to change with right input")
	}
}

func TestSpaceBrakeReducesVelocity(t *testing.T) {
	bs := newTestBootstrap(rand.New(rand.NewSource(1)))
	ctx := ecs.GetWorldContext(bs.World)
	ctx.CurrentLayer = ecs.LayerSpace
	ecs.SetWorldContext(bs.World, ctx)
	SetPlayerInput(bs.World, bs.Player, "up")
	for i := 0; i < 60; i++ {
		bs.Scheduler.Update(1.0/20.0, bs.World)
	}
	v0, _ := ecs.Get[components.Velocity](bs.World, bs.Player)
	SetPlayerInput(bs.World, bs.Player, "down")
	for i := 0; i < 30; i++ {
		bs.Scheduler.Update(1.0/20.0, bs.World)
	}
	v1, _ := ecs.Get[components.Velocity](bs.World, bs.Player)
	if abs(v1.VX)+abs(v1.VY) >= abs(v0.VX)+abs(v0.VY) {
		t.Fatal("expected brake to reduce speed")
	}
}
