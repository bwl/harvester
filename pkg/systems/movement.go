package systems

import (
	"bubbleRouge/pkg/components"
	"bubbleRouge/pkg/ecs"
)

type Movement struct{}

func (Movement) Update(dt float64, w *ecs.World) {
	ecs.View2Of[components.Position, components.Velocity](w).Each(func(t ecs.Tuple2[components.Position, components.Velocity]) {
		t.A.X += t.B.VX * dt
		t.A.Y += t.B.VY * dt
		ecs.Add(w, t.E, *t.A)
	})
}
