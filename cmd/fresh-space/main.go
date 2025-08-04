package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"harvester/pkg/engine"
)

// FreshSpaceModel is a completely minimal space scene with spaceship physics
type FreshSpaceModel struct {
	world  *ecs.World
	player ecs.Entity
	width  int
	height int
	// Spaceship state
	angle    float64                // rotation angle in radians
	velocity struct{ x, y float64 } // momentum
	thrust   bool                   // is thrusting
	braking  bool                   // is braking
}

func NewFreshSpaceModel() *FreshSpaceModel {
	// Create ECS world with player
	bs := engine.New(nil)

	// Get the player position and move it to center
	if pos, ok := ecs.Get[components.Position](bs.World, bs.Player); ok {
		pos.X = 50.0
		pos.Y = 25.0
		ecs.Add(bs.World, bs.Player, pos)
	}

	// Reset velocity to prevent drift
	if vel, ok := ecs.Get[components.Velocity](bs.World, bs.Player); ok {
		vel.VX = 0.0
		vel.VY = 0.0
		ecs.Add(bs.World, bs.Player, vel)
	}

	return &FreshSpaceModel{
		world:  bs.World,
		player: bs.Player,
		width:  80,
		height: 24,
		angle:  -math.Pi / 2, // facing north initially (up)
	}
}

func (m *FreshSpaceModel) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg { return t }) // ~60 FPS
}

func (m *FreshSpaceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case time.Time:
		// Physics update
		m.updatePhysics()
		return m, tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg { return t })

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "w":
			m.thrust = true // thrust forward
		case "s":
			m.braking = true // brake
		case "a":
			m.angle -= 0.1 // turn left
		case "d":
			m.angle += 0.1 // turn right
		case "c":
			// Center player and stop
			if pos, ok := ecs.Get[components.Position](m.world, m.player); ok {
				pos.X = float64(m.width / 2)
				pos.Y = float64(m.height / 2)
				ecs.Add(m.world, m.player, pos)
			}
			m.velocity.x = 0
			m.velocity.y = 0
		}
		return m, nil
	}

	m.thrust = false  // only thrust when key is held
	m.braking = false // only brake when key is held
	return m, nil
}

func (m *FreshSpaceModel) updatePhysics() {
	const dt = 0.016           // ~60 FPS
	const thrustPower = 30.0   // Much lower thrust power
	const baseFriction = 0.98  // Base friction (less drag)
	const brakeFriction = 0.92 // Brake friction (more drag than base)
	const maxSpeed = 80.0      // Maximum speed limit

	// Apply thrust
	if m.thrust {
		thrustX := math.Cos(m.angle) * thrustPower * dt
		thrustY := math.Sin(m.angle) * thrustPower * dt
		m.velocity.x += thrustX
		m.velocity.y += thrustY
	}

	// Apply friction - brake friction is STRONGER (lower number) than base friction
	frictionToUse := baseFriction
	if m.braking {
		frictionToUse = brakeFriction
	}

	m.velocity.x *= frictionToUse
	m.velocity.y *= frictionToUse

	// Stop if very slow
	speed := math.Sqrt(m.velocity.x*m.velocity.x + m.velocity.y*m.velocity.y)
	if speed < 1.0 {
		m.velocity.x = 0
		m.velocity.y = 0
	}

	// Limit maximum speed
	speed = math.Sqrt(m.velocity.x*m.velocity.x + m.velocity.y*m.velocity.y)
	if speed > maxSpeed {
		m.velocity.x = (m.velocity.x / speed) * maxSpeed
		m.velocity.y = (m.velocity.y / speed) * maxSpeed
	}

	// Update position based on velocity
	if pos, ok := ecs.Get[components.Position](m.world, m.player); ok {
		newPos := pos // Copy the position
		newPos.X += m.velocity.x * dt
		newPos.Y += m.velocity.y * dt

		// Wrap around screen edges
		if newPos.X < 0 {
			newPos.X = float64(m.width)
		} else if newPos.X > float64(m.width) {
			newPos.X = 0
		}
		if newPos.Y < 0 {
			newPos.Y = float64(m.height)
		} else if newPos.Y > float64(m.height) {
			newPos.Y = 0
		}

		ecs.Add(m.world, m.player, newPos)
	}
}

func (m *FreshSpaceModel) View() string {
	// Create a simple star field with the player
	lines := make([]string, m.height)

	// Fill with spaces
	for i := range lines {
		lines[i] = strings.Repeat(" ", m.width)
	}

	// Add some stars
	for y := 0; y < m.height; y += 3 {
		for x := 0; x < m.width; x += 7 {
			if y < len(lines) && x < len(lines[y]) {
				lineBytes := []byte(lines[y])
				lineBytes[x] = '.'
				lines[y] = string(lineBytes)
			}
		}
	}

	// Add the player with directional indicator
	if pos, ok := ecs.Get[components.Position](m.world, m.player); ok {
		px, py := int(pos.X), int(pos.Y)
		if py >= 0 && py < len(lines) && px >= 0 && px < len(lines[py]) {
			lineBytes := []byte(lines[py])

			// Choose ship glyph based on angle using Unicode arrows
			var shipGlyph string
			// Normalize angle to 0-2π
			normalizedAngle := math.Mod(m.angle+2*math.Pi, 2*math.Pi)

			if normalizedAngle < math.Pi/8 || normalizedAngle >= 15*math.Pi/8 {
				shipGlyph = "🢂" // facing right
			} else if normalizedAngle < 3*math.Pi/8 {
				shipGlyph = "🢆" // facing down-right
			} else if normalizedAngle < 5*math.Pi/8 {
				shipGlyph = "🢃" // facing down
			} else if normalizedAngle < 7*math.Pi/8 {
				shipGlyph = "🢇" // facing down-left
			} else if normalizedAngle < 9*math.Pi/8 {
				shipGlyph = "🢀" // facing left
			} else if normalizedAngle < 11*math.Pi/8 {
				shipGlyph = "🢄" // facing up-left
			} else if normalizedAngle < 13*math.Pi/8 {
				shipGlyph = "🢁" // facing up
			} else {
				shipGlyph = "🢅" // facing up-right
			}

			// Replace the character at this position with the ship glyph
			line := string(lineBytes)
			if px > 0 {
				line = line[:px] + shipGlyph + line[px+1:]
			} else {
				line = shipGlyph + line[1:]
			}
			lines[py] = line
		}
	}

	return strings.Join(lines, "\n")
}

func main() {
	fmt.Println("Spaceship Flight Simulator")
	fmt.Println("Controls:")
	fmt.Println("  W - Thrust forward")
	fmt.Println("  A/D - Turn left/right")
	fmt.Println("  S - Brake")
	fmt.Println("  C - Center and stop")
	fmt.Println("  Q - Quit")
	fmt.Println()
	fmt.Println("Ship direction: 🢁 🢃 🢀 🢂 🢄 🢅 🢆 🢇")
	fmt.Println("Press any key to start...")
	fmt.Scanln()

	model := NewFreshSpaceModel()
	program := tea.NewProgram(model, tea.WithAltScreen())

	if err := program.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
