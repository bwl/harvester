package systems

import (
	"bubbleRouge/pkg/components"
	"bubbleRouge/pkg/ecs"
)

type Drawable struct {
	X, Y  int
	Glyph rune
}

type Render struct{ Output []Drawable }

func (r *Render) Update(dt float64, w *ecs.World) {
	out := make([]Drawable, 0, 128)
	ecs.View2Of[components.Position, components.Renderable](w).Each(func(t ecs.Tuple2[components.Position, components.Renderable]) {
		out = append(out, Drawable{X: int(t.A.X), Y: int(t.A.Y), Glyph: t.B.Glyph})
	})
	r.Output = out
}
