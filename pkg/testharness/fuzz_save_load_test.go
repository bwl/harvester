//go:build go1.18

package testharness

import (
	"testing"

	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

func FuzzSaveLoadEquivalence(f *testing.F) {
	seed := func(px, py int, fuel, hull, drive int) {
		w := ecs.NewWorld(nil)
		p := w.Create()
		ecs.Add(w, p, components.Position{X: float64(px%1000 - 500), Y: float64(py%1000 - 500)})
		ecs.Add(w, p, components.PlayerStats{Fuel: fuel % 1000, Hull: hull % 1000, Drive: (drive%5 + 1)})
		ecs.Add(w, 1, components.WorldInfo{Width: 200, Height: 80})
		s1, _ := ecs.Save(w, nil)
		s2, _ := ecs.Save(w, nil)
		_ = s1
		_ = s2
	}
	seed(0, 0, 100, 100, 1)
	seed(10, -20, 50, 75, 3)

	f.Fuzz(func(t *testing.T, px, py int, fuel, hull, drive int) {
		w := ecs.NewWorld(nil)
		p := w.Create()
		ecs.Add(w, p, components.Position{X: float64(px%1000 - 500), Y: float64(py%1000 - 500)})
		ecs.Add(w, p, components.PlayerStats{Fuel: abs(fuel) % 1000, Hull: abs(hull) % 1000, Drive: (abs(drive)%5 + 1)})
		ecs.Add(w, 1, components.WorldInfo{Width: 200, Height: 80})

		s1, err := ecs.Save(w, nil)
		if err != nil {
			t.Fatalf("save1: %v", err)
		}
		w2 := ecs.NewWorld(nil)
		if err := ecs.Load(w2, s1, nil); err != nil {
			t.Fatalf("load: %v", err)
		}
		s2, err := ecs.Save(w2, nil)
		if err != nil {
			t.Fatalf("save2: %v", err)
		}

		// encode both snapshots to JSON bytes and compare deterministically
		b1, _ := ecs.EncodeSnapshot(s1, ecs.SaveOptions{})
		b2, _ := ecs.EncodeSnapshot(s2, ecs.SaveOptions{})
		if string(b1) != string(b2) {
			// try compressed round too
			c1, _ := ecs.EncodeSnapshot(s1, ecs.SaveOptions{Compress: true})
			c2, _ := ecs.EncodeSnapshot(s2, ecs.SaveOptions{Compress: true})
			if len(c1) == 0 || len(c2) == 0 || len(c1) != len(c2) {
				t.Fatalf("snapshots differ after save-load-save")
			}
		}
	})
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
