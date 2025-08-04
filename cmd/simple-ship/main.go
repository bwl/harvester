package main

import (
	"fmt"
	"math"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/internal/ui"
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"harvester/pkg/engine"
	"harvester/pkg/rendering"
	"harvester/pkg/systems"
)

// SimpleShip - integrates real ECS world with simple physics
type SimpleShip struct {
	x, y    float64 // position
	vx, vy  float64 // velocity
	angle   float64 // rotation
	width   int
	height  int
	thrust  bool
	braking bool
	// Camera system
	cameraX, cameraY float64 // camera world position
	// Renderer
	renderer *rendering.ViewRenderer
	// ECS world for expanse tiles
	world     *ecs.World
	player    ecs.Entity
	render    *systems.Render
	scheduler *ecs.SchedulerWithContext
}

// Planet represents a planet in space
type Planet struct {
	x, y  float64
	glyph string
	name  string
}

func NewSimpleShip() *SimpleShip {
	// Create ECS world with real game engine
	bs := engine.New(nil)

	s := &SimpleShip{
		x: 40, y: 12, // center of screen
		angle: -math.Pi / 2, // facing up
		width: 80, height: 24,
		// ECS components
		world:     bs.World,
		player:    bs.Player,
		render:    bs.Render,
		scheduler: bs.Scheduler,
	}

	// Initialize camera to center on player
	s.cameraX = s.x - float64(s.width)/2
	s.cameraY = s.y - float64(s.height)/2

	// Initialize renderer
	s.renderer = rendering.NewViewRenderer(s.width, s.height)

	// Set player position in ECS to match physics position
	if pos, ok := ecs.Get[components.Position](s.world, s.player); ok {
		pos.X = s.x
		pos.Y = s.y
		ecs.Add(s.world, s.player, pos)
	}

	// Make sure we're in space layer for expanse tiles
	ctx := ecs.GetWorldContext(s.world)
	ctx.CurrentLayer = ecs.LayerSpace
	ecs.SetWorldContext(s.world, ctx)

	return s
}

func (s *SimpleShip) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg { return t })
}

func (s *SimpleShip) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		// Update renderer dimensions
		s.renderer = rendering.NewViewRenderer(s.width, s.height)
		return s, nil

	case time.Time:
		s.updatePhysics()
		// Run ECS scheduler to update all systems (including expansion!)
		s.scheduler.Update(0.016, s.world)
		// Reset thrust/brake after physics update
		s.thrust = false
		s.braking = false
		return s, tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg { return t })

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return s, tea.Quit
		case "w":
			s.thrust = true
		case "s":
			s.braking = true
		case "a":
			s.angle -= 0.15
		case "d":
			s.angle += 0.15
		case "c":
			// Reset to world origin and stop
			s.x = 0
			s.y = 0
			s.vx = 0
			s.vy = 0
		case "p":
			// Debug: print position info and render system status
			fmt.Printf("Player: (%.1f, %.1f), Camera: (%.1f, %.1f), Velocity: (%.1f, %.1f)\n",
				s.x, s.y, s.cameraX, s.cameraY, s.vx, s.vy)
			fmt.Printf("Render output count: %d\n", len(s.render.Output))
			if len(s.render.Output) > 0 {
				fmt.Printf("First few render items: ")
				for i, d := range s.render.Output {
					if i >= 3 {
						break
					}
					fmt.Printf("(%d,%d:'%c') ", d.X, d.Y, d.Glyph)
				}
				fmt.Printf("\n")
			}
			// Check world context
			ctx := ecs.GetWorldContext(s.world)
			fmt.Printf("Current Layer: %v, Entity count: %d\n", ctx.CurrentLayer, s.world.EntityCount())
		}
		return s, nil
	}

	return s, nil
}

func (s *SimpleShip) updatePhysics() {
	const dt = 0.016
	const thrustPower = 30.0
	const friction = 0.999     // Much less general resistance
	const brakeFriction = 0.92 // Keep brake the same
	const maxSpeed = 80.0

	// Apply thrust
	if s.thrust {
		thrustX := math.Cos(s.angle) * thrustPower * dt
		thrustY := math.Sin(s.angle) * thrustPower * dt
		s.vx += thrustX
		s.vy += thrustY
	}

	// Apply friction
	if s.braking {
		s.vx *= brakeFriction
		s.vy *= brakeFriction
	} else {
		s.vx *= friction
		s.vy *= friction
	}

	// Speed limit
	speed := math.Sqrt(s.vx*s.vx + s.vy*s.vy)
	if speed > maxSpeed {
		s.vx = (s.vx / speed) * maxSpeed
		s.vy = (s.vy / speed) * maxSpeed
	}

	// Stop if very slow
	if speed < 0.1 {
		s.vx = 0
		s.vy = 0
	}

	// Update position
	s.x += s.vx * dt
	s.y += s.vy * dt

	// Update camera to follow player (keep player centered)
	s.cameraX = s.x - float64(s.width)/2
	s.cameraY = s.y - float64(s.height)/2

	// No wrap around - let ship fly in infinite space
	// Camera will follow, showing different areas
}

func (s *SimpleShip) buildGlyphMatrix() [][]rendering.Glyph {
	// Create glyph matrix
	glyphs := make([][]rendering.Glyph, s.height)
	for i := range glyphs {
		glyphs[i] = make([]rendering.Glyph, s.width)
		// Fill with space glyphs
		for j := range glyphs[i] {
			glyphs[i][j] = rendering.Glyph{
				Char:       ' ',
				Foreground: rendering.Color{R: 255, G: 255, B: 255},
				Background: rendering.Color{R: 0, G: 0, B: 0},
				Alpha:      1.0,
			}
		}
	}

	// Add stars based on camera position for infinite star field
	starOffsetX := int(s.cameraX) % 7
	starOffsetY := int(s.cameraY) % 3
	for y := -starOffsetY; y < s.height; y += 3 {
		for x := -starOffsetX; x < s.width; x += 7 {
			screenX, screenY := x, y
			if screenY >= 0 && screenY < s.height && screenX >= 0 && screenX < s.width {
				glyphs[screenY][screenX] = rendering.Glyph{
					Char:       '.',
					Foreground: rendering.Color{R: 100, G: 100, B: 100}, // Dim gray stars
					Background: rendering.Color{R: 0, G: 0, B: 0},
					Alpha:      1.0,
				}
			}
		}
	}

	// Add planets (in world coordinates, convert to screen coordinates)
	planets := []Planet{
		{x: 20, y: 8, glyph: "🌍", name: "Earth"},
		{x: 60, y: 18, glyph: "🪐", name: "Saturn"},
		{x: 15, y: 20, glyph: "🔴", name: "Mars"},
		{x: 70, y: 5, glyph: "🌕", name: "Moon"},
	}

	for _, planet := range planets {
		// Convert world coordinates to screen coordinates
		screenX := int(planet.x - s.cameraX)
		screenY := int(planet.y - s.cameraY)

		if screenY >= 0 && screenY < s.height && screenX >= 0 && screenX < s.width {
			// Get planet color based on type
			var color rendering.Color
			switch planet.name {
			case "Earth":
				color = rendering.Color{R: 0, G: 100, B: 255} // Blue
			case "Mars":
				color = rendering.Color{R: 255, G: 50, B: 50} // Red
			case "Saturn":
				color = rendering.Color{R: 200, G: 180, B: 100} // Golden
			case "Moon":
				color = rendering.Color{R: 220, G: 220, B: 220} // Light gray
			default:
				color = rendering.Color{R: 255, G: 255, B: 255} // White
			}

			glyphs[screenY][screenX] = rendering.Glyph{
				Char:       []rune(planet.glyph)[0], // Convert emoji to rune
				Foreground: color,
				Background: rendering.Color{R: 0, G: 0, B: 0},
				Alpha:      1.0,
			}
		}
	}

	// Add ship (always centered on screen)
	shipScreenX := s.width / 2
	shipScreenY := s.height / 2

	if shipScreenY >= 0 && shipScreenY < s.height && shipScreenX >= 0 && shipScreenX < s.width {
		shipGlyph := s.getShipGlyph()
		glyphs[shipScreenY][shipScreenX] = rendering.Glyph{
			Char:       []rune(shipGlyph)[0],
			Foreground: rendering.Color{R: 255, G: 255, B: 0}, // Bright yellow ship
			Background: rendering.Color{R: 0, G: 0, B: 0},
			Style:      rendering.StyleBold,
			Alpha:      1.0,
		}
	}

	return glyphs
}

// expanseContent - similar to the original UI renderer

// GetExpanseContent returns renderable content using real ECS system
func (s *SimpleShip) GetExpanseContent() rendering.LayerContent {
	// Sync ECS player position with physics position
	if pos, ok := ecs.Get[components.Position](s.world, s.player); ok {
		pos.X = s.x
		pos.Y = s.y
		ecs.Add(s.world, s.player, pos)
	}

	// Update camera in ECS to match our camera system
	if cam, ok := ecs.Get[components.Camera](s.world, s.player); ok {
		cam.X = int(s.cameraX)
		cam.Y = int(s.cameraY)
		ecs.Add(s.world, s.player, cam)
	}

	// Use the real buildGameGlyphs from the game system
	glyphs := s.buildGameGlyphs(s.width, s.height)
	if glyphs == nil {
		return nil
	}
	return ui.NewExpanseContent(glyphs)
}

func (s *SimpleShip) buildGameGlyphs(w, h int) [][]rendering.Glyph {
	if w <= 0 || h <= 0 {
		return nil
	}
	// Use unified render system for all drawables (expanse tiles!)
	s.render.Update(0, s.world)
	cam, _ := ecs.Get[components.Camera](s.world, s.player)
	mx0, my0 := cam.X, cam.Y

	glyphs := make([][]rendering.Glyph, h)
	for y := 0; y < h; y++ {
		row := make([]rendering.Glyph, w)
		for x := 0; x < w; x++ {
			row[x] = rendering.Glyph{Char: '.'}
		}
		glyphs[y] = row
	}

	// Render all ECS entities (this includes expanse tiles!)
	rendered := 0
	for _, d := range s.render.Output {
		x := d.X - mx0
		y := d.Y - my0
		if x >= 0 && y >= 0 && x < w && y < h {
			glyphs[y][x] = rendering.Glyph{Char: rune(d.Glyph)}
			rendered++
		}
	}

	// Debug: print coordinate translation (only when tiles exist but none rendered)
	if len(s.render.Output) > 1 && rendered == 1 {
		fmt.Printf("[DEBUG] Camera: (%d,%d), Screen: %dx%d, Rendered %d/%d items\n",
			mx0, my0, w, h, rendered, len(s.render.Output))
		for i, d := range s.render.Output {
			if i >= 5 {
				break
			}
			screenX, screenY := d.X-mx0, d.Y-my0
			visible := screenX >= 0 && screenY >= 0 && screenX < w && screenY < h
			fmt.Printf("  Item %d: World(%d,%d) -> Screen(%d,%d) '%c' visible=%v\n",
				i, d.X, d.Y, screenX, screenY, d.Glyph, visible)
		}
	}

	return glyphs
}

func (s *SimpleShip) View() string {
	// Register our expanse content with the global renderer
	content := s.GetExpanseContent()
	if content != nil {
		s.renderer.RegisterContent(content)
	}

	// Use the global renderer to generate the final string output
	return s.renderer.Render()
}

func (s *SimpleShip) getShipGlyph() string {
	normalizedAngle := math.Mod(s.angle+2*math.Pi, 2*math.Pi)

	if normalizedAngle < math.Pi/8 || normalizedAngle >= 15*math.Pi/8 {
		return "🢂"
	} else if normalizedAngle < 3*math.Pi/8 {
		return "🢆"
	} else if normalizedAngle < 5*math.Pi/8 {
		return "🢃"
	} else if normalizedAngle < 7*math.Pi/8 {
		return "🢇"
	} else if normalizedAngle < 9*math.Pi/8 {
		return "🢀"
	} else if normalizedAngle < 11*math.Pi/8 {
		return "🢄"
	} else if normalizedAngle < 13*math.Pi/8 {
		return "🢁"
	} else {
		return "🢅"
	}
}

func main() {
	fmt.Println("Simple Ship Test - Camera Version")
	fmt.Println("Testing camera system...")
	testCamera()
	fmt.Println("\nStarting game...")
	fmt.Println("W=thrust, S=brake, A/D=turn, C=reset to origin, P=debug, Q=quit")
	fmt.Println("Press any key to start...")
	fmt.Scanln()

	ship := NewSimpleShip()
	program := tea.NewProgram(ship, tea.WithAltScreen())

	if err := program.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
