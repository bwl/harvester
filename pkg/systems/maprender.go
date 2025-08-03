package systems

import (
	"bubbleRouge/pkg/components"
	"bubbleRouge/pkg/ecs"
)

type MapRender struct{ Output []Drawable }

func (m *MapRender) Update(dt float64, w *ecs.World) {
	out := m.Output[:0]
	ecs.View2Of[components.Position, components.Tile](w).Each(func(t ecs.Tuple2[components.Position, components.Tile]) {
		out = append(out, Drawable{X: int(t.A.X), Y: int(t.A.Y), Glyph: t.B.Glyph})
	})
	m.Output = out
}
