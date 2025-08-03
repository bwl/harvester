package systems

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

type Combat struct{}

func (Combat) Update(dt float64, w *ecs.World) {
	// Demo: if Left pressed and adjacent enemy, apply damage
	ecs.View2Of[components.Input, components.Position](w).Each(func(t ecs.Tuple2[components.Input, components.Position]) {
		if !t.A.Left {
			return
		}
		px, py := int(t.B.X)-1, int(t.B.Y)
		ecs.View2Of[components.Position, components.Health](w).Each(func(u ecs.Tuple2[components.Position, components.Health]) {
			if int(u.A.X) == px && int(u.A.Y) == py {
				u.B.HP -= 10
				if u.B.HP < 0 {
					u.B.HP = 0
				}
				ecs.Add(w, u.E, *u.B)
			}
		})
	})
}
