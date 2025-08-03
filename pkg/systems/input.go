package systems

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

type InputSystem struct{}

type Control struct{ Entity ecs.Entity }

type EnterPlanet struct{}

func (InputSystem) Update(dt float64, w *ecs.World) {
	const thrustRamp = 40.0
	const thrustDecay = 20.0
	const turnRate = 2.5
	ctx := ecs.GetWorldContext(w)
	isSpace := ctx.CurrentLayer == ecs.LayerSpace
	if isSpace {
		ecs.View1Of[components.Input](w).Each(func(e ecs.Entity, in *components.Input) {
			spr, _ := ecs.Get[components.SpaceFlightSprings](w, e)
			if in.Up {
				spr.Thrust.Target = 100
			} else {
				spr.Thrust.Target = 0
			}
			if in.Left {
				spr.Angle.Target -= turnRate
			} else if in.Right {
				spr.Angle.Target += turnRate
			}
			if in.Down {
				spr.VelX.Target, spr.VelY.Target = 0, 0
			}
			ecs.Add(w, e, spr)
		})
		return
	}
	// non-space fallback
	ecs.View2Of[components.Input, components.Acceleration](w).Each(func(t ecs.Tuple2[components.Input, components.Acceleration]) {
		ax, ay := 0.0, 0.0
		if t.A.Left {
			ax -= thrustRamp
		}
		if t.A.Right {
			ax += thrustRamp
		}
		if t.A.Up {
			ay -= thrustRamp
		}
		if t.A.Down {
			ay += thrustRamp
		}
		t.B.AX, t.B.AY = ax, ay
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
