package systems

import (
	"bubbleRouge/pkg/components"
	"bubbleRouge/pkg/ecs"
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
	ecs.Add(w, c.Target, cam)
}
