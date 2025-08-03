package testharness

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

func TestHarvest(t *testing.T) {
	c := NewController(Options{Width: 40, Height: 20})
	e := c.World.Create()
	ecs.Add(c.World, e, components.Position{})
	ecs.Add(c.World, e, components.Resource{Kind: "ore", Amount: 1})
	ecs.Add(c.World, c.Player, components.Action{Harvest: true})
	c.Tick(1, 1)
	inv, _ := ecs.Get[components.Inventory](c.World, c.Player)
	require.Equal(t, 1, inv.Items["ore"])
}

func TestSaveLoad(t *testing.T) {
	c := NewController(Options{})
	c.InjectKey("right")
	c.Tick(5, 1)
	s1, err := c.Snapshot()
	require.NoError(t, err)
	var snap1 map[string]any
	require.NoError(t, json.Unmarshal(s1, &snap1))
	b, err := ecs.Save(c.World, nil)
	require.NoError(t, err)
	// mutate
	c.InjectKey("left")
	c.Tick(3, 1)
	require.NoError(t, ecs.Load(c.World, b, nil))
	s2, _ := c.Snapshot()
	var snap2 map[string]any
	_ = json.Unmarshal(s2, &snap2)
	require.Equal(t, snap1["player"], snap2["player"])
}
