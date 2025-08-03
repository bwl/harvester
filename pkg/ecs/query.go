package ecs

type View1[A any] struct{ as *store[A] }

type View2[A any, B any] struct {
	as *store[A]
	bs *store[B]
}

type Tuple1[A any] struct {
	E Entity
	A *A
}

type Tuple2[A any, B any] struct {
	E Entity
	A *A
	B *B
}

type View3[A any, B any, C any] struct {
	as *store[A]
	bs *store[B]
	cs *store[C]
}

type Tuple3[A any, B any, C any] struct {
	E Entity
	A *A
	B *B
	C *C
}

func View1Of[A any](w *World) View1[A] { return View1[A]{as: storeOf[A](w)} }
func View2Of[A any, B any](w *World) View2[A, B] {
	return View2[A, B]{as: storeOf[A](w), bs: storeOf[B](w)}
}
func View3Of[A any, B any, C any](w *World) View3[A, B, C] {
	return View3[A, B, C]{as: storeOf[A](w), bs: storeOf[B](w), cs: storeOf[C](w)}
}

func (v View1[A]) Each(fn func(e Entity, a *A)) { v.as.ForEach(func(e Entity, a *A) { fn(e, a) }) }

func (v View2[A, B]) Each(fn func(t Tuple2[A, B])) {
	v.as.ForEach(func(e Entity, a *A) {
		if b, ok := v.bs.Get(e); ok {
			bb := b
			fn(Tuple2[A, B]{E: e, A: a, B: &bb})
		}
	})
}

func (v View3[A, B, C]) Each(fn func(t Tuple3[A, B, C])) {
	v.as.ForEach(func(e Entity, a *A) {
		if b, ok := v.bs.Get(e); ok {
			if c, ok := v.cs.Get(e); ok {
				bb := b
				cc := c
				fn(Tuple3[A, B, C]{E: e, A: a, B: &bb, C: &cc})
			}
		}
	})
}
