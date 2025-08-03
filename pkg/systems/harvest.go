package systems

import (
	"bubbleRouge/pkg/components"
	"bubbleRouge/pkg/ecs"
)

type Harvest struct{}

func (Harvest) Update(dt float64, w *ecs.World) {
	// Triggered by Action.Harvest
	ecs.View2Of[components.Action, components.Position](w).Each(func(t ecs.Tuple2[components.Action, components.Position]) {
		if !t.A.Harvest {
			return
		}
		// find resource at same position
		target := ecs.Entity(0)
		var res components.Resource
		ecs.View2Of[components.Position, components.Resource](w).Each(func(u ecs.Tuple2[components.Position, components.Resource]) {
			if int(u.A.X) == int(t.B.X) && int(u.A.Y) == int(t.B.Y) {
				target = u.E
				res = *u.B
			}
		})
		if target == 0 {
			return
		}
		inv, _ := ecs.Get[components.Inventory](w, t.E)
		if inv.Items == nil {
			inv.Items = make(map[string]int)
		}
		inv.Items[res.Kind] += res.Amount
		ecs.Add(w, t.E, inv)
		ecs.Remove[components.Resource](w, target)
	})
}
