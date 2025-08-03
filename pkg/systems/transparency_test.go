package systems

import (
	"math/rand"
	"testing"

	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

func TestTransparencyComponent(t *testing.T) {
	// Test transparency component creation
	trans := components.NewTransparency(0.5)
	if trans.Alpha != 0.5 {
		t.Errorf("Expected alpha 0.5, got %f", trans.Alpha)
	}
	if trans.BlendMode != components.BlendNormal {
		t.Errorf("Expected BlendNormal, got %v", trans.BlendMode)
	}

	// Test with blend mode
	transBlend := components.NewTransparencyWithBlend(0.3, components.BlendAdditive)
	if transBlend.Alpha != 0.3 {
		t.Errorf("Expected alpha 0.3, got %f", transBlend.Alpha)
	}
	if transBlend.BlendMode != components.BlendAdditive {
		t.Errorf("Expected BlendAdditive, got %v", transBlend.BlendMode)
	}
}

func TestMapRenderWithTransparency(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	world := ecs.NewWorld(r)
	
	// Create entity with transparency
	entity := world.Create()
	ecs.Add(world, entity, components.Position{X: 5, Y: 5})
	ecs.Add(world, entity, components.Tile{
		Glyph: '#',
		Type:  components.TileForest,
	})
	ecs.Add(world, entity, components.Transparency{
		Alpha:     0.7,
		BlendMode: components.BlendNormal,
	})

	// Create map renderer
	mapRender := &MapRender{}
	mapRender.Update(0.0, world)

	// Verify alpha was applied
	if len(mapRender.Output) != 1 {
		t.Fatalf("Expected 1 drawable, got %d", len(mapRender.Output))
	}

	drawable := mapRender.Output[0]
	if drawable.Alpha != 0.7 {
		t.Errorf("Expected alpha 0.7, got %f", drawable.Alpha)
	}
	if drawable.BlendMode != components.BlendNormal {
		t.Errorf("Expected BlendNormal, got %v", drawable.BlendMode)
	}
}

func TestMapRenderWithoutTransparency(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	world := ecs.NewWorld(r)
	
	// Create entity without transparency
	entity := world.Create()
	ecs.Add(world, entity, components.Position{X: 5, Y: 5})
	ecs.Add(world, entity, components.Tile{
		Glyph: '#',
		Type:  components.TileForest,
	})

	// Create map renderer
	mapRender := &MapRender{}
	mapRender.Update(0.0, world)

	// Verify default alpha (opaque)
	if len(mapRender.Output) != 1 {
		t.Fatalf("Expected 1 drawable, got %d", len(mapRender.Output))
	}

	drawable := mapRender.Output[0]
	if drawable.Alpha != 1.0 {
		t.Errorf("Expected alpha 1.0 (opaque), got %f", drawable.Alpha)
	}
	if drawable.BlendMode != components.BlendNormal {
		t.Errorf("Expected BlendNormal, got %v", drawable.BlendMode)
	}
}

func TestFadeEffect(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	world := ecs.NewWorld(r)
	
	// Create entity with fade effect
	entity := world.Create()
	ecs.Add(world, entity, FadeEffect{
		StartAlpha: 1.0,
		EndAlpha:   0.0,
		Duration:   2.0,
		Elapsed:    0.0,
	})

	// Create fade system
	fadeSystem := &FadeEffectSystem{}
	
	// Update partway through animation
	fadeSystem.Update(1.0, world) // 1 second elapsed, halfway
	
	// Check transparency was added
	if trans, ok := ecs.Get[components.Transparency](world, entity); ok {
		expectedAlpha := 0.5 // Halfway between 1.0 and 0.0
		if trans.Alpha != expectedAlpha {
			t.Errorf("Expected alpha %f, got %f", expectedAlpha, trans.Alpha)
		}
	} else {
		t.Error("Expected transparency component to be added")
	}
	
	// Update to completion
	fadeSystem.Update(1.0, world) // Another 1 second, total 2 seconds
	
	// Check final alpha
	if trans, ok := ecs.Get[components.Transparency](world, entity); ok {
		if trans.Alpha != 0.0 {
			t.Errorf("Expected final alpha 0.0, got %f", trans.Alpha)
		}
	} else {
		t.Error("Expected transparency component to remain")
	}
	
	// Check fade effect was removed
	if _, ok := ecs.Get[FadeEffect](world, entity); ok {
		t.Error("Expected fade effect to be removed after completion")
	}
}

func TestHelperFunctions(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	world := ecs.NewWorld(r)
	
	// Test CreateFog
	fogEntity := CreateFog(world, 10, 10, 0.8)
	
	if pos, ok := ecs.Get[components.Position](world, fogEntity); ok {
		if pos.X != 10 || pos.Y != 10 {
			t.Errorf("Expected position (10, 10), got (%f, %f)", pos.X, pos.Y)
		}
	} else {
		t.Error("Expected fog entity to have position")
	}
	
	if trans, ok := ecs.Get[components.Transparency](world, fogEntity); ok {
		expectedAlpha := 0.3 * 0.8 // intensity affects opacity
		if trans.Alpha != expectedAlpha {
			t.Errorf("Expected alpha %f, got %f", expectedAlpha, trans.Alpha)
		}
	} else {
		t.Error("Expected fog entity to have transparency")
	}
	
	// Test StartFadeOut
	StartFadeOut(world, fogEntity, 1.0)
	
	if fade, ok := ecs.Get[FadeEffect](world, fogEntity); ok {
		if fade.Duration != 1.0 {
			t.Errorf("Expected fade duration 1.0, got %f", fade.Duration)
		}
		if fade.EndAlpha != 0.0 {
			t.Errorf("Expected end alpha 0.0, got %f", fade.EndAlpha)
		}
	} else {
		t.Error("Expected fade effect to be added")
	}
}