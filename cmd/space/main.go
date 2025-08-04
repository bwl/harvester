package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/internal/ui"
	"harvester/pkg/components"
	"harvester/pkg/debug"
	"harvester/pkg/ecs"
)

// SimpleSpaceView directly loads spaceScreen without going through start menu
type SimpleSpaceView struct {
	spaceScreen    *ui.SpaceScreen
	model          *ui.Model
	w, h           int
	needsCentering bool
}

func NewSimpleSpaceView() *SimpleSpaceView {
	// Create a new model (this initializes the ECS world and systems)
	model := ui.NewModel(nil)

	// Create the space screen with this model
	spaceScreen := ui.NewSpaceScreen(&model)

	return &SimpleSpaceView{
		spaceScreen:    spaceScreen,
		model:          &model,
		needsCentering: true,
	}
}

func (s *SimpleSpaceView) Init() tea.Cmd {
	return s.spaceScreen.Init()
}

func (s *SimpleSpaceView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Center player on first update after everything is initialized
	if s.needsCentering {
		fmt.Printf("*** CENTERING PLAYER NOW ***\n")
		s.centerPlayer()
		s.needsCentering = false
		fmt.Printf("*** CENTERING COMPLETE ***\n")
	}

	// Handle window sizing
	if wm, ok := msg.(tea.WindowSizeMsg); ok {
		s.w, s.h = wm.Width, wm.Height
		s.spaceScreen.SetDimensions(s.w, s.h)
		fmt.Printf("*** WINDOW RESIZE: %dx%d ***\n", s.w, s.h)
	}

	// Handle quit keys
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "ctrl+c", "q":
			return s, tea.Quit
		case "c":
			// Manual centering on 'c' key press
			fmt.Printf("*** MANUAL CENTERING TRIGGERED ***\n")
			s.centerPlayer()
			return s, nil
		}
	}

	// Forward to space screen
	newModel, cmd := s.spaceScreen.Update(msg)
	if spaceScreen, ok := newModel.(*ui.SpaceScreen); ok {
		s.spaceScreen = spaceScreen
	}
	return s, cmd
}

func (s *SimpleSpaceView) centerPlayer() {
	world := s.model.World()

	// Move player to center of screen and reset all physics
	newX, newY := 50.0, 25.0
	debug.Info("space", fmt.Sprintf("Centering player at %.1f,%.1f and resetting physics", newX, newY))

	playerCount := 0
	view := ecs.View2Of[components.Player, components.Position](world)
	view.Each(func(t ecs.Tuple2[components.Player, components.Position]) {
		playerCount++
		debug.Info("space", fmt.Sprintf("Before: Player entity %d at (%.1f, %.1f)", t.E, t.B.X, t.B.Y))

		// Set new position
		t.B.X = newX
		t.B.Y = newY

		// Reset velocity to zero
		if vel, ok := ecs.Get[components.Velocity](world, t.E); ok {
			vel.VX = 0.0
			vel.VY = 0.0
			ecs.Add(world, t.E, vel)
			debug.Info("space", "Reset velocity to zero")
		}

		// Reset acceleration to zero
		if accel, ok := ecs.Get[components.Acceleration](world, t.E); ok {
			accel.AX = 0.0
			accel.AY = 0.0
			ecs.Add(world, t.E, accel)
			debug.Info("space", "Reset acceleration to zero")
		}

		// Reset springs if they exist
		if springs, ok := ecs.Get[components.SpaceFlightSprings](world, t.E); ok {
			// Reset spring positions and velocities to zero
			springs.VelX.Pos = 0.0
			springs.VelX.Vel = 0.0
			springs.VelX.Target = 0.0
			springs.VelY.Pos = 0.0
			springs.VelY.Vel = 0.0
			springs.VelY.Target = 0.0
			springs.Thrust.Pos = 0.0
			springs.Thrust.Vel = 0.0
			springs.Thrust.Target = 0.0
			ecs.Add(world, t.E, springs)
			debug.Info("space", "Reset spring physics to zero")
		}

		debug.Info("space", fmt.Sprintf("After: Player entity %d moved to (%.1f, %.1f) with physics reset", t.E, t.B.X, t.B.Y))
	})

	if playerCount == 0 {
		debug.Warn("space", "No player entities found!")
	} else {
		debug.Info("space", fmt.Sprintf("Centered %d player entities with physics reset", playerCount))
	}
}

func (s *SimpleSpaceView) View() string {
	return s.spaceScreen.View()
}

func main() {
	// Initialize debugging system
	debug.Info("space", "Starting Space Screen Test")
	debug.Info("space", "Press 'q' or Ctrl+C to quit")

	// Create simple space view that bypasses start screen
	spaceView := NewSimpleSpaceView()

	// Launch the application with space view directly
	program := tea.NewProgram(spaceView, tea.WithAltScreen())
	if err := program.Start(); err != nil {
		debug.Errorf("space", "Application failed to start: %v", err)
		fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
		os.Exit(1)
	}

	debug.Info("space", "Space screen test shutting down")
}
