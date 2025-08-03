package testharness

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"harvester/pkg/engine"
	"math/rand"
	"testing"
)

func TestSpaceMovementTick(t *testing.T) {
	bs := engine.New(rand.New(rand.NewSource(1)))
	w := bs.World
	p := bs.Player
	ctx := ecs.GetWorldContext(w)
	ctx.CurrentLayer = ecs.LayerSpace
	ecs.SetWorldContext(w, ctx)
	ecs.Add(w, p, components.Input{Up: true})
	for i := 0; i < 40; i++ {
		bs.Scheduler.Update(0.05, w)
	}
	pos, _ := ecs.Get[components.Position](w, p)
	if pos.X <= 0 {
		t.Fatalf("expected movement, got %v", pos.X)
	}
}

func TestSaveLoad_EngineComponents(t *testing.T) {
	bs := engine.New(rand.New(rand.NewSource(1)))
	w := bs.World
	s, err := ecs.Save(w, nil)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ecs.EncodeSnapshot(s, ecs.SaveOptions{Compress: true})
	if err != nil {
		t.Fatal(err)
	}
	s2, err := ecs.DecodeSnapshot(b, ecs.SaveOptions{Compress: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := ecs.Load(w, s2, nil); err != nil {
		t.Fatal(err)
	}
}
