package ecs

import (
	"encoding/json"
	"math/rand"
	"reflect"

	"harvester/pkg/components"
)

type Snapshot struct {
	Version    int                                   `json:"version"`
	Seed       int64                                 `json:"seed"`
	Next       Entity                                `json:"next"`
	Free       []Entity                              `json:"free"`
	Components map[string]map[Entity]json.RawMessage `json:"components"`
}

func Save(w *World, enc func(v any) ([]byte, error)) (*Snapshot, error) {
	if enc == nil {
		enc = json.Marshal
	}
	s := &Snapshot{Components: make(map[string]map[Entity]json.RawMessage)}
	s.Version = currentSnapshotVersion()
	// entity allocator state
	w.mu.RLock()
	s.Seed = w.seed
	s.Next = w.next
	if len(w.free) > 0 {
		cp := make([]Entity, len(w.free))
		copy(cp, w.free)
		s.Free = cp
	}
	w.mu.RUnlock()
	// Persist known baseline component stores explicitly for type safety
	s.Components[typeName[components.Position]()] = dumpStore(enc, storeOf[components.Position](w))
	s.Components[typeName[components.Velocity]()] = dumpStore(enc, storeOf[components.Velocity](w))
	s.Components[typeName[components.Camera]()] = dumpStore(enc, storeOf[components.Camera](w))
	s.Components[typeName[components.PlayerStats]()] = dumpStore(enc, storeOf[components.PlayerStats](w))
	s.Components[typeName[components.WorldInfo]()] = dumpStore(enc, storeOf[components.WorldInfo](w))
	s.Components[typeName[components.Input]()] = dumpStore(enc, storeOf[components.Input](w))
	s.Components[typeName[components.Inventory]()] = dumpStore(enc, storeOf[components.Inventory](w))
	s.Components[typeName[components.Resource]()] = dumpStore(enc, storeOf[components.Resource](w))
	s.Components[typeName[components.Tile]()] = dumpStore(enc, storeOf[components.Tile](w))
	s.Components[typeName[components.Renderable]()] = dumpStore(enc, storeOf[components.Renderable](w))
	s.Components[typeName[components.Health]()] = dumpStore(enc, storeOf[components.Health](w))
	// persist WorldContext and surface-related systems' ad hoc components
	s.Components[typeName[WorldContext]()] = dumpStore(enc, storeOf[WorldContext](w))
	return s, nil
}

func Load(w *World, s *Snapshot, dec func(data []byte, v any) error) error {
	if dec == nil {
		dec = json.Unmarshal
	}
	if err := maybeMigrateSnapshot(s); err != nil {
		return err
	}
	// restore allocator deterministically
	w.mu.Lock()
	if s.Version >= 1 && s.Seed != 0 {
		w.seed = s.Seed
		w.rng = rand.New(rand.NewSource(w.seed))
	}
	w.next = s.Next
	w.free = nil
	if len(s.Free) > 0 {
		w.free = make([]Entity, len(s.Free))
		copy(w.free, s.Free)
	}
	w.mu.Unlock()
	loadStore(dec, storeOf[components.Position](w), s.Components[typeName[components.Position]()])
	loadStore(dec, storeOf[components.Velocity](w), s.Components[typeName[components.Velocity]()])
	loadStore(dec, storeOf[components.Camera](w), s.Components[typeName[components.Camera]()])
	loadStore(dec, storeOf[components.PlayerStats](w), s.Components[typeName[components.PlayerStats]()])
	loadStore(dec, storeOf[components.WorldInfo](w), s.Components[typeName[components.WorldInfo]()])
	loadStore(dec, storeOf[components.Input](w), s.Components[typeName[components.Input]()])
	loadStore(dec, storeOf[components.Inventory](w), s.Components[typeName[components.Inventory]()])
	loadStore(dec, storeOf[components.Resource](w), s.Components[typeName[components.Resource]()])
	loadStore(dec, storeOf[components.Tile](w), s.Components[typeName[components.Tile]()])
	loadStore(dec, storeOf[components.Renderable](w), s.Components[typeName[components.Renderable]()])
	loadStore(dec, storeOf[components.Health](w), s.Components[typeName[components.Health]()])
	loadStore(dec, storeOf[WorldContext](w), s.Components[typeName[WorldContext]()])
	return nil
}

func dumpStore[T any](enc func(v any) ([]byte, error), st *store[T]) map[Entity]json.RawMessage {
	m := make(map[Entity]json.RawMessage)
	st.ForEach(func(e Entity, t *T) {
		b, _ := enc(t)
		m[e] = b
	})
	return m
}

func loadStore[T any](dec func([]byte, any) error, st *store[T], data map[Entity]json.RawMessage) {
	if data == nil {
		return
	}
	// clear existing for deterministic restore
	st.mu.Lock()
	st.data = make(map[Entity]T)
	st.index = make(map[Entity]struct{})
	st.mu.Unlock()
	for e, raw := range data {
		var v T
		if dec(raw, &v) == nil {
			// post-unmarshal fixups for known types
			switch any(&v).(type) {
			case *components.Inventory:
				iv := any(&v).(*components.Inventory)
				iv.Ensure()
			}
			st.Add(e, v)
		}
	}
}

func typeName[T any]() string { return reflect.TypeOf((*T)(nil)).Elem().String() }

var snapshotMigrations = map[int]func(*Snapshot) error{}

func currentSnapshotVersion() int { return 1 }

func maybeMigrateSnapshot(s *Snapshot) error {
	v := s.Version
	for v < currentSnapshotVersion() {
		mig, ok := snapshotMigrations[v]
		if !ok {
			break
		}
		if err := mig(s); err != nil {
			return err
		}
		v++
		s.Version = v
	}
	return nil
}
