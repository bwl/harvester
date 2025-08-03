package ecs

import (
	"math/rand"
	"reflect"
	"sync"
)

type World struct {
	mu     sync.RWMutex
	next   Entity
	free   []Entity
	stores map[reflect.Type]any
	rng    *rand.Rand
	seed   int64
	saveMu sync.Mutex
}

func NewWorld(r *rand.Rand) *World {
	if r == nil {
		r = rand.New(rand.NewSource(1))
	}
	return &World{stores: make(map[reflect.Type]any), rng: r, seed: 1}
}

func RandFromSeed(seed int64) *rand.Rand { return rand.New(rand.NewSource(seed)) }

func (w *World) Create() Entity {
	w.mu.Lock()
	defer w.mu.Unlock()
	if n := len(w.free); n > 0 {
		e := w.free[n-1]
		w.free = w.free[:n-1]
		return e
	}
	w.next++
	return w.next
}

func (w *World) Destroy(e Entity) {
	w.mu.Lock()
	for _, st := range w.stores {
		removeFromStore(st, e)
	}
	w.free = append(w.free, e)
	w.mu.Unlock()
}

func storeOf[T any](w *World) *store[T] {
	t := reflect.TypeOf((*T)(nil)).Elem()
	st, ok := w.stores[t]
	if !ok {
		ss := newStore[T]()
		w.stores[t] = ss
		return ss
	}
	return st.(*store[T])
}

func Add[T any](w *World, e Entity, c T)      { storeOf[T](w).Add(e, c) }
func Get[T any](w *World, e Entity) (T, bool) { return storeOf[T](w).Get(e) }
func Remove[T any](w *World, e Entity)        { storeOf[T](w).Remove(e) }

func removeFromStore(st any, e Entity) {
	switch s := st.(type) {
	case *store[int]:
		s.Remove(e)
	default:
		// Using reflection to call Remove on generic store
		rv := reflect.ValueOf(st)
		m := rv.MethodByName("Remove")
		if m.IsValid() {
			m.Call([]reflect.Value{reflect.ValueOf(e)})
		}
	}
}
