package engine

import (
	"testing"
)

func TestNewInitializesWorld(t *testing.T) {
	gs := New(42)
	if gs.Map == nil || gs.Map.Width <= 0 || gs.Map.Height <= 0 {
		t.Fatalf("invalid map: %+v", gs.Map)
	}
	if gs.Player.Pos.X < 0 || gs.Player.Pos.Y < 0 {
		t.Fatalf("player pos invalid: %+v", gs.Player.Pos)
	}
	found := 0
	for y := 0; y < gs.Map.Height; y++ {
		for x := 0; x < gs.Map.Width; x++ {
			if gs.Map.Tiles[y][x].Kind == Galaxy { found++ }
		}
	}
	if found == 0 {
		t.Fatalf("expected some initial galaxies, got %d", found)
	}
}

func TestExpansionOnMove(t *testing.T) {
	gs := New(1)
	w0, h0 := gs.Map.Width, gs.Map.Height
	gs.Step(ActMoveRight)
	if gs.Map.Width <= w0 || gs.Map.Height <= h0 {
		t.Fatalf("expected expansion, got %dx%d -> %dx%d", w0, h0, gs.Map.Width, gs.Map.Height)
	}
	// player should remain in bounds after shift
	if gs.Player.Pos.X < 0 || gs.Player.Pos.Y < 0 || gs.Player.Pos.X >= gs.Map.Width || gs.Player.Pos.Y >= gs.Map.Height {
		t.Fatalf("player out of bounds after expansion: %+v size:%dx%d", gs.Player.Pos, gs.Map.Width, gs.Map.Height)
	}
}

func TestHarvestRemovesGalaxy(t *testing.T) {
	gs := New(99)
	// place a galaxy at player
	p := gs.Player.Pos
	gs.Map.Tiles[p.Y][p.X].Kind = Galaxy
	gs.Step(ActHarvest)
	if gs.Map.Tiles[p.Y][p.X].Kind == Galaxy {
		t.Fatalf("expected galaxy harvested at %v", p)
	}
}

func TestDeterministicWithSeed(t *testing.T) {
	gs1 := New(123)
	gs2 := New(123)
	// perform same actions
	for i := 0; i < 5; i++ { gs1.Step(ActMoveRight) }
	for i := 0; i < 5; i++ { gs2.Step(ActMoveRight) }
	if gs1.Map.Width != gs2.Map.Width || gs1.Map.Height != gs2.Map.Height || gs1.Player.Pos != gs2.Player.Pos {
		t.Fatalf("determinism failed: gs1(%dx%d,%v) gs2(%dx%d,%v)", gs1.Map.Width, gs1.Map.Height, gs1.Player.Pos, gs2.Map.Width, gs2.Map.Height, gs2.Player.Pos)
	}
}
