package main

import (
	"fmt"
	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

func testCamera() {
	ship := NewSimpleShip()

	fmt.Printf("Initial: Player(%.1f, %.1f), Camera(%.1f, %.1f)\n",
		ship.x, ship.y, ship.cameraX, ship.cameraY)

	// Debug ECS system
	ctx := ecs.GetWorldContext(ship.world)
	fmt.Printf("ECS World - Layer: %v, Entity count: %d\n", ctx.CurrentLayer, ship.world.EntityCount())

	// Test render system
	ship.render.Update(0, ship.world)
	fmt.Printf("Render output count: %d\n", len(ship.render.Output))
	if len(ship.render.Output) > 0 {
		fmt.Printf("First few render items: ")
		for i, d := range ship.render.Output {
			if i >= 5 {
				break
			}
			fmt.Printf("(%d,%d:'%c') ", d.X, d.Y, d.Glyph)
		}
		fmt.Printf("\n")
	} else {
		fmt.Printf("No render output - checking if expansion system works...\n")
		// Try to trigger expansion
		if pos, ok := ecs.Get[components.Position](ship.world, ship.player); ok {
			fmt.Printf("Player position in ECS: (%.1f, %.1f)\n", pos.X, pos.Y)
		}
	}

	// Move ship with thrust
	ship.thrust = true
	for i := 0; i < 10; i++ {
		ship.updatePhysics()
		if i%5 == 0 {
			fmt.Printf("Step %d: Player(%.1f, %.1f), Camera(%.1f, %.1f)\n",
				i, ship.x, ship.y, ship.cameraX, ship.cameraY)
		}
	}

	// Test render system after movement and run scheduler
	ship.scheduler.Update(0.016, ship.world)
	ship.render.Update(0, ship.world)
	fmt.Printf("After movement and scheduler - Render output count: %d\n", len(ship.render.Output))
	if len(ship.render.Output) > 1 {
		fmt.Printf("SUCCESS! Found expanse tiles: ")
		for i, d := range ship.render.Output {
			if i >= 10 {
				break
			}
			fmt.Printf("(%d,%d:'%c') ", d.X, d.Y, d.Glyph)
		}
		fmt.Printf("\n")

		// Test glyph matrix generation
		fmt.Printf("Testing glyph matrix generation...\n")
		glyphs := ship.buildGameGlyphs(ship.width, ship.height)
		if glyphs != nil {
			nonSpaceCount := 0
			for y := 0; y < len(glyphs); y++ {
				for x := 0; x < len(glyphs[y]); x++ {
					if glyphs[y][x].Char != '.' && glyphs[y][x].Char != ' ' && glyphs[y][x].Char != 0 {
						nonSpaceCount++
						if nonSpaceCount <= 5 {
							fmt.Printf("  Found non-space glyph at (%d,%d): '%c'\n", x, y, glyphs[y][x].Char)
						}
					}
				}
			}
			fmt.Printf("Total non-space glyphs in matrix: %d\n", nonSpaceCount)
		}
	}
}

// Removed main function to avoid conflicts
