package systems

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

// FadeEffect component for animated transparency changes
type FadeEffect struct {
	StartAlpha float64
	EndAlpha   float64
	Duration   float64
	Elapsed    float64
}

// FadeEffectSystem handles animated fade effects
type FadeEffectSystem struct{}

func (f *FadeEffectSystem) Update(dt float64, w *ecs.World) {
	ecs.View1Of[FadeEffect](w).Each(func(e ecs.Entity, fade *FadeEffect) {
		fade.Elapsed += dt
		progress := fade.Elapsed / fade.Duration

		if progress >= 1.0 {
			// Animation complete
			finalAlpha := fade.EndAlpha
			ecs.Add(w, e, components.Transparency{
				Alpha:     finalAlpha,
				BlendMode: components.BlendNormal,
			})
			ecs.Remove[FadeEffect](w, e)
		} else {
			// Interpolate alpha
			currentAlpha := fade.StartAlpha + (fade.EndAlpha-fade.StartAlpha)*progress
			ecs.Add(w, e, components.Transparency{
				Alpha:     currentAlpha,
				BlendMode: components.BlendNormal,
			})
			ecs.Add(w, e, *fade)
		}
	})
}

// DecaySystem for fading corpses and debris
type DecaySystem struct{}

type DecayTimer struct {
	Duration float64
	Elapsed  float64
}

func (d *DecaySystem) Update(dt float64, w *ecs.World) {
	ecs.View1Of[DecayTimer](w).Each(func(e ecs.Entity, timer *DecayTimer) {
		timer.Elapsed += dt

		// Fade out over time
		fadeProgress := timer.Elapsed / timer.Duration
		alpha := 1.0 - fadeProgress

		if alpha <= 0 {
			// Remove completely decayed entity
			w.Destroy(e)
		} else {
			// Update transparency
			ecs.Add(w, e, components.Transparency{
				Alpha:     alpha,
				BlendMode: components.BlendNormal,
			})
			ecs.Add(w, e, *timer)
		}
	})
}

// Helper functions for creating transparent effects

// CreateFog creates a fog entity with transparency
func CreateFog(world *ecs.World, x, y float64, intensity float64) ecs.Entity {
	fogEntity := world.Create()

	ecs.Add(world, fogEntity, components.Position{X: x, Y: y})
	ecs.Add(world, fogEntity, components.Tile{
		Glyph: '░',
		Type:  components.TileForest, // Reuse existing tile type for now
	})
	ecs.Add(world, fogEntity, components.Transparency{
		Alpha:     0.3 * intensity, // Intensity affects opacity
		BlendMode: components.BlendNormal,
	})

	return fogEntity
}

// CreateSmoke creates a smoke entity with additive blending
func CreateSmoke(world *ecs.World, x, y float64, intensity float64) ecs.Entity {
	smokeEntity := world.Create()

	ecs.Add(world, smokeEntity, components.Position{X: x, Y: y})
	ecs.Add(world, smokeEntity, components.Renderable{
		Glyph:    '▒',
		TileType: components.TileForest, // Reuse existing tile type for now
	})
	ecs.Add(world, smokeEntity, components.Transparency{
		Alpha:     0.4 * intensity,          // Intensity affects opacity
		BlendMode: components.BlendAdditive, // Smoke adds to background
	})

	return smokeEntity
}

// StartFadeOut starts a fade out effect on an entity
func StartFadeOut(world *ecs.World, entity ecs.Entity, duration float64) {
	// Get current alpha or default to 1.0
	startAlpha := 1.0
	if trans, ok := ecs.Get[components.Transparency](world, entity); ok {
		startAlpha = trans.Alpha
	}

	ecs.Add(world, entity, FadeEffect{
		StartAlpha: startAlpha,
		EndAlpha:   0.0,
		Duration:   duration,
		Elapsed:    0.0,
	})
}

// StartFadeIn starts a fade in effect on an entity
func StartFadeIn(world *ecs.World, entity ecs.Entity, duration float64) {
	// Get current alpha or default to 0.0
	startAlpha := 0.0
	if trans, ok := ecs.Get[components.Transparency](world, entity); ok {
		startAlpha = trans.Alpha
	}

	ecs.Add(world, entity, FadeEffect{
		StartAlpha: startAlpha,
		EndAlpha:   1.0,
		Duration:   duration,
		Elapsed:    0.0,
	})
}
