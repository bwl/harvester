package systems

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

type Render struct{ Output []Drawable }

func (r *Render) Update(dt float64, w *ecs.World) {
	out := make([]Drawable, 0, 128)
	ecs.View2Of[components.Position, components.Renderable](w).Each(func(t ecs.Tuple2[components.Position, components.Renderable]) {
		out = append(out, Drawable{X: int(t.A.X), Y: int(t.A.Y), Glyph: t.B.Glyph})
	})
	r.Output = out
}
