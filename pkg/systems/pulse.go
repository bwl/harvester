package systems

import (
	"github.com/charmbracelet/harmonica"
	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

type PulseSystem struct{ spring harmonica.Spring }

func (p *PulseSystem) ensure() {
	if p.spring == (harmonica.Spring{}) {
		p.spring = harmonica.NewSpring(harmonica.FPS(60), 5.0, 0.2)
	}
}

func (p *PulseSystem) Update(dt float64, w *ecs.World) {
	p.ensure()
	ecs.View1Of[components.PulseSpring](w).Each(func(e ecs.Entity, ps *components.PulseSpring) {
		ps.Pos, ps.Vel = p.spring.Update(ps.Pos, ps.Vel, ps.Target)
		ecs.Add(w, e, *ps)
	})
}
