package systems

import (
	"harvester/pkg/ecs"
	"harvester/pkg/timing"
)

type Tick struct{}

func (Tick) Update(dt float64, w *ecs.World) {
	// Update global timer
	timing.UpdateGlobalTimer()
}
