package testharness

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

func TestConcurrentSaveLoad(t *testing.T) {
	w := ecs.NewWorld(nil)
	p := w.Create()
	ecs.Add(w, p, components.Position{X: 0, Y: 0})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				s, err := ecs.Save(w, nil)
				require.NoError(t, err)
				w2 := ecs.NewWorld(nil)
				require.NoError(t, ecs.Load(w2, s, nil))
			}
		}()
	}
	wg.Wait()
}
