package systems

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

type CameraSystem struct{ Target ecs.Entity }

func (c *CameraSystem) Update(dt float64, w *ecs.World) {
	pos, ok := ecs.Get[components.Position](w, c.Target)
	if !ok {
		return
	}
	cam, _ := ecs.Get[components.Camera](w, c.Target)
	cam.X = int(pos.X) - cam.Width/2
	cam.Y = int(pos.Y) - cam.Height/2
	wi, _ := ecs.Get[components.WorldInfo](w, 1)
	if cam.X < 0 {
		cam.X = 0
	}
	if cam.Y < 0 {
		cam.Y = 0
	}
	maxX := wi.Width - cam.Width
	maxY := wi.Height - cam.Height
	if cam.Width > wi.Width {
		cam.X = 0
	} else if cam.X > maxX {
		cam.X = maxX
	}
	if cam.Height > wi.Height {
		cam.Y = 0
	} else if cam.Y > maxY {
		cam.Y = maxY
	}
	ecs.Add(w, c.Target, cam)
}
