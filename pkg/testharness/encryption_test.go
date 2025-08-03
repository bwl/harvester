package testharness

import (
	"testing"

	"bubbleRouge/pkg/components"
	"bubbleRouge/pkg/ecs"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecode_WithPassword(t *testing.T) {
	w := ecs.NewWorld(nil)
	p := w.Create()
	ecs.Add(w, p, components.Position{X: 1, Y: 2})
	s, err := ecs.Save(w, nil)
	require.NoError(t, err)
	b, err := ecs.EncodeSnapshot(s, ecs.SaveOptions{Compress: true, Password: "pw"})
	require.NoError(t, err)
	s2, err := ecs.DecodeSnapshot(b, ecs.SaveOptions{Compress: true, Password: "pw"})
	require.NoError(t, err)
	w2 := ecs.NewWorld(nil)
	require.NoError(t, ecs.Load(w2, s2, nil))
	p2, ok := ecs.Get[components.Position](w2, p)
	require.True(t, ok)
	require.Equal(t, 1.0, p2.X)
	require.Equal(t, 2.0, p2.Y)
}
