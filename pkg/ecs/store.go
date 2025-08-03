package ecs

import "sync"

type store[T any] struct {
	mu    sync.RWMutex
	data  map[Entity]T
	index map[Entity]struct{}
}

func newStore[T any]() *store[T] {
	return &store[T]{
		data:  make(map[Entity]T),
		index: make(map[Entity]struct{}),
	}
}

func (s *store[T]) Add(e Entity, c T) {
	s.mu.Lock()
	s.data[e] = c
	s.index[e] = struct{}{}
	s.mu.Unlock()
}

func (s *store[T]) Get(e Entity) (T, bool) {
	s.mu.RLock()
	v, ok := s.data[e]
	s.mu.RUnlock()
	return v, ok
}

func (s *store[T]) Remove(e Entity) {
	s.mu.Lock()
	delete(s.data, e)
	delete(s.index, e)
	s.mu.Unlock()
}

func (s *store[T]) Has(e Entity) bool {
	s.mu.RLock()
	_, ok := s.index[e]
	s.mu.RUnlock()
	return ok
}

func (s *store[T]) ForEach(f func(Entity, *T)) {
	s.mu.RLock()
	for e := range s.index {
		v := s.data[e]
		vv := v
		s.mu.RUnlock()
		f(e, &vv)
		s.mu.RLock()
	}
	s.mu.RUnlock()
}
