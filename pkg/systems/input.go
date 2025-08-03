package systems

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

type InputSystem struct{}

type Control struct{ Entity ecs.Entity }

type EnterPlanet struct{}

func (InputSystem) Update(dt float64, w *ecs.World) {
	// translate Input -> Velocity for all layers
	ecs.View2Of[components.Input, components.Velocity](w).Each(func(t ecs.Tuple2[components.Input, components.Velocity]) {
		vx, vy := 0.0, 0.0
		if t.A.Left {
			vx = -1
		}
		if t.A.Right {
			vx = 1
		}
		if t.A.Up {
			vy = -1
		}
		if t.A.Down {
			vy = 1
		}
		t.B.VX, t.B.VY = vx, vy
		ecs.Add(w, t.E, *t.B)
	})
}

func SetPlayerInput(w *ecs.World, e ecs.Entity, dir string) {
	in, _ := ecs.Get[components.Input](w, e)
	switch dir {
	case "left":
		in.Left, in.Right, in.Up, in.Down = true, false, false, false
	case "right":
		in.Left, in.Right, in.Up, in.Down = false, true, false, false
	case "up":
		in.Left, in.Right, in.Up, in.Down = false, false, true, false
	case "down":
		in.Left, in.Right, in.Up, in.Down = false, false, false, true
	case "enter":
		in.Left, in.Right, in.Up, in.Down = false, false, false, false
		ecs.Add(w, e, EnterPlanet{})
	case "clear":
		// no-op retain last state
	default:
		in = components.Input{}
	}
	ecs.Add(w, e, in)
}
