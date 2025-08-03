package testharness

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMovementAndCamera(t *testing.T) {
	c := NewController(Options{Width: 40, Height: 20})
	c.InjectKey("right")
	c.Tick(10, 1)
	snap, err := c.Snapshot()
	require.NoError(t, err)
	require.Contains(t, string(snap), "\"x\": 10")
	require.Contains(t, string(snap), "\"camera\":")
}
