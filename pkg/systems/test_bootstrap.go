package systems

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"math/rand"
)

type testBootstrap struct {
	World     *ecs.World
	Scheduler *ecs.SchedulerWithContext
	Player    ecs.Entity
}

func newTestBootstrap(r *rand.Rand) testBootstrap {
	w := ecs.NewWorld(r)
	p := w.Create()
	ecs.Add(w, p, components.Player{})
	ecs.Add(w, p, components.Position{})
	ecs.Add(w, p, components.Input{})
	ecs.Add(w, p, components.Velocity{})
	ecs.Add(w, p, components.Acceleration{})
	ecs.Add(w, p, components.SpaceFlightSprings{})
	reg := ecs.SystemRegistry{
		UniversalSystems: []ecs.System{InputSystem{}, Tick{}},
		SpaceSystems:     []ecs.System{SpaceMovement{}},
	}
	s := ecs.NewSchedulerWithContext(reg)
	ecs.Add(w, 1, components.WorldInfo{Width: 200, Height: 80})
	ecs.SetWorldContext(w, ecs.WorldContext{CurrentLayer: ecs.LayerSpace})
	return testBootstrap{World: w, Scheduler: s, Player: p}
}
