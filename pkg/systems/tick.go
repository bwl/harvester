package systems

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

type Tick struct{}

func (Tick) Update(dt float64, w *ecs.World) {
	wi, ok := ecs.Get[components.WorldInfo](w, 1)
	if !ok {
		ecs.Add(w, 1, components.WorldInfo{})
		wi, _ = ecs.Get[components.WorldInfo](w, 1)
	}
	wi.Tick++
	ecs.Add(w, 1, wi)
}
