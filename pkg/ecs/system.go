package ecs

type System interface {
	Update(dt float64, w *World)
}

type Scheduler struct {
	order []System
}

func NewScheduler(order ...System) *Scheduler {
	return &Scheduler{order: order}
}

func (s *Scheduler) Update(dt float64, w *World) {
	for _, sys := range s.order {
		sys.Update(dt, w)
	}
}
