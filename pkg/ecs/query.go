package ecs

type View2[A any, B any] struct {
	as *store[A]
	bs *store[B]
}

type Tuple2[A any, B any] struct {
	E Entity
	A *A
	B *B
}

func View2Of[A any, B any](w *World) View2[A, B] {
	return View2[A, B]{as: storeOf[A](w), bs: storeOf[B](w)}
}

func (v View2[A, B]) Each(fn func(t Tuple2[A, B])) {
	v.as.ForEach(func(e Entity, a *A) {
		if b, ok := v.bs.Get(e); ok {
			bb := b
			fn(Tuple2[A, B]{E: e, A: a, B: &bb})
		}
	})
}
