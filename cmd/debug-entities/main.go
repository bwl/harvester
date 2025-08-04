package main

import (
	"fmt"
	"harvester/internal/ui"
	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

func main() {
	fmt.Println("=== Entity Debug Test ===")

	// Create a model like our space script does
	model := ui.NewModel(nil)
	world := model.World()

	fmt.Println("World created, examining entities...")

	// Check world info
	if wi, ok := ecs.Get[components.WorldInfo](world, 1); ok {
		fmt.Printf("World info: %dx%d\n", wi.Width, wi.Height)
	} else {
		fmt.Println("No world info found!")
	}

	// Look for player entities
	fmt.Println("\n=== Searching for Player entities ===")
	playerView := ecs.View1Of[components.Player](world)
	playerCount := 0
	playerView.Each(func(e ecs.Entity, p *components.Player) {
		playerCount++
		fmt.Printf("Found Player entity: %d\n", e)

		// Check if it has Position
		if pos, ok := ecs.Get[components.Position](world, e); ok {
			fmt.Printf("  Position: (%.1f, %.1f)\n", pos.X, pos.Y)
		} else {
			fmt.Printf("  No Position component!\n")
		}

		// Check if it has Camera
		if cam, ok := ecs.Get[components.Camera](world, e); ok {
			fmt.Printf("  Camera: (%d, %d) size %dx%d\n", cam.X, cam.Y, cam.Width, cam.Height)
		} else {
			fmt.Printf("  No Camera component!\n")
		}

		// Check if it has Renderable
		if rend, ok := ecs.Get[components.Renderable](world, e); ok {
			fmt.Printf("  Renderable: '%c'\n", rend.Glyph)
		} else {
			fmt.Printf("  No Renderable component!\n")
		}
	})

	fmt.Printf("\nTotal Player entities found: %d\n", playerCount)

	// Look for all Position entities
	fmt.Println("\n=== All Position entities ===")
	posView := ecs.View1Of[components.Position](world)
	posCount := 0
	posView.Each(func(e ecs.Entity, pos *components.Position) {
		posCount++
		fmt.Printf("Entity %d: Position (%.1f, %.1f)\n", e, pos.X, pos.Y)
	})

	fmt.Printf("\nTotal Position entities found: %d\n", posCount)

	// Now try to modify player position
	fmt.Println("\n=== Attempting to center player ===")
	view := ecs.View2Of[components.Player, components.Position](world)
	view.Each(func(t ecs.Tuple2[components.Player, components.Position]) {
		fmt.Printf("Before: Entity %d at (%.1f, %.1f)\n", t.E, t.B.X, t.B.Y)
		t.B.X = 100.0
		t.B.Y = 40.0
		fmt.Printf("After:  Entity %d at (%.1f, %.1f)\n", t.E, t.B.X, t.B.Y)
	})

	// Verify the change stuck
	fmt.Println("\n=== Verification ===")
	view.Each(func(t ecs.Tuple2[components.Player, components.Position]) {
		fmt.Printf("Final: Entity %d at (%.1f, %.1f)\n", t.E, t.B.X, t.B.Y)
	})
}
