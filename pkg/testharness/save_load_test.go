package testharness

import (
	"encoding/json"
	"testing"

	"bubbleRouge/pkg/components"
	"bubbleRouge/pkg/ecs"
	"github.com/stretchr/testify/require"
)

func roundtrip(t *testing.T, w *ecs.World) *ecs.World {
	s, err := ecs.Save(w, nil)
	require.NoError(t, err)
	w2 := ecs.NewWorld(nil)
	require.NoError(t, ecs.Load(w2, s, nil))
	return w2
}

func TestSaveLoad_Position_Velocity(t *testing.T) {
	w := ecs.NewWorld(nil)
	e := w.Create()
	ecs.Add(w, e, components.Position{X: 3, Y: 4})
	ecs.Add(w, e, components.Velocity{VX: 1, VY: -2})
	w2 := roundtrip(t, w)
	p, ok := ecs.Get[components.Position](w2, e)
	require.True(t, ok)
	v, ok := ecs.Get[components.Velocity](w2, e)
	require.True(t, ok)
	require.Equal(t, 3.0, p.X)
	require.Equal(t, 4.0, p.Y)
	require.Equal(t, 1.0, v.VX)
	require.Equal(t, -2.0, v.VY)
}

func TestSaveLoad_Camera_PlayerStats_WorldInfo(t *testing.T) {
	w := ecs.NewWorld(nil)
	player := w.Create()
	ecs.Add(w, player, components.Camera{X: 5, Y: 6, Width: 40, Height: 20})
	ecs.Add(w, player, components.PlayerStats{Fuel: 7, Hull: 8, Drive: 2})
	ecs.Add(w, 1, components.WorldInfo{Tick: 42, Width: 200, Height: 80})
	w2 := roundtrip(t, w)
	cam, _ := ecs.Get[components.Camera](w2, player)
	ps, _ := ecs.Get[components.PlayerStats](w2, player)
	wi, _ := ecs.Get[components.WorldInfo](w2, 1)
	require.Equal(t, 5, cam.X)
	require.Equal(t, 6, cam.Y)
	require.Equal(t, 40, cam.Width)
	require.Equal(t, 20, cam.Height)
	require.Equal(t, 7, ps.Fuel)
	require.Equal(t, 8, ps.Hull)
	require.Equal(t, 2, ps.Drive)
	require.Equal(t, int64(42), wi.Tick)
	require.Equal(t, 200, wi.Width)
	require.Equal(t, 80, wi.Height)
}

func TestSaveLoad_Inventory(t *testing.T) {
	w := ecs.NewWorld(nil)
	p := w.Create()
	inv := components.Inventory{}
	inv.Ensure()
	inv.Items["ore"] = 3
	ecs.Add(w, p, inv)
	w2 := roundtrip(t, w)
	inv2, ok := ecs.Get[components.Inventory](w2, p)
	require.True(t, ok)
	require.NotNil(t, inv2.Items)
	require.Equal(t, 3, inv2.Items["ore"])
}

func TestSaveLoad_Tile_Renderable_Health_Resource(t *testing.T) {
	w := ecs.NewWorld(nil)
	e := w.Create()
	ecs.Add(w, e, components.Tile{Glyph: '*'})
	ecs.Add(w, e, components.Renderable{Glyph: '@'})
	ecs.Add(w, e, components.Health{HP: 5, Max: 10})
	ecs.Add(w, e, components.Resource{Kind: "ore", Amount: 9})
	w2 := roundtrip(t, w)
	tile, _ := ecs.Get[components.Tile](w2, e)
	r, _ := ecs.Get[components.Renderable](w2, e)
	h, _ := ecs.Get[components.Health](w2, e)
	res, _ := ecs.Get[components.Resource](w2, e)
	require.Equal(t, '*', tile.Glyph)
	require.Equal(t, '@', r.Glyph)
	require.Equal(t, 5, h.HP)
	require.Equal(t, 10, h.Max)
	require.Equal(t, "ore", res.Kind)
	require.Equal(t, 9, res.Amount)
}

func TestController_SaveLoad_EndToEnd(t *testing.T) {
	c := NewController(Options{Width: 40, Height: 20})
	c.InjectKey("right")
	c.Tick(5, 1)
	s1, _ := c.Snapshot()
	var a, b map[string]any
	_ = json.Unmarshal(s1, &a)
	snap, err := ecs.Save(c.World, nil)
	require.NoError(t, err)
	c.InjectKey("left")
	c.Tick(3, 1)
	require.NoError(t, ecs.Load(c.World, snap, nil))
	s2, _ := c.Snapshot()
	_ = json.Unmarshal(s2, &b)
	require.Equal(t, a["player"], b["player"])
}
