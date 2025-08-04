package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/internal/ui"
	"harvester/pkg/debug"
)

// BasicSpaceView - no centering, just load space screen
type BasicSpaceView struct {
	spaceScreen *ui.SpaceScreen
	model       *ui.Model
	w, h        int
}

func NewBasicSpaceView() *BasicSpaceView {
	// Create a new model (this initializes the ECS world and systems)
	model := ui.NewModel(nil)

	// Create the space screen with this model
	spaceScreen := ui.NewSpaceScreen(&model)

	return &BasicSpaceView{
		spaceScreen: spaceScreen,
		model:       &model,
	}
}

func (s *BasicSpaceView) Init() tea.Cmd {
	return s.spaceScreen.Init()
}

func (s *BasicSpaceView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle window sizing
	if wm, ok := msg.(tea.WindowSizeMsg); ok {
		s.w, s.h = wm.Width, wm.Height
		s.spaceScreen.SetDimensions(s.w, s.h)
	}

	// Handle quit keys
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "ctrl+c", "q":
			return s, tea.Quit
		}
	}

	// Forward to space screen
	newModel, cmd := s.spaceScreen.Update(msg)
	if spaceScreen, ok := newModel.(*ui.SpaceScreen); ok {
		s.spaceScreen = spaceScreen
	}
	return s, cmd
}

func (s *BasicSpaceView) View() string {
	return s.spaceScreen.View()
}

func main() {
	debug.Info("basic-space", "Starting Basic Space Screen Test")
	debug.Info("basic-space", "No centering - player should be at original position")

	// Create basic space view
	spaceView := NewBasicSpaceView()

	// Launch the application
	program := tea.NewProgram(spaceView, tea.WithAltScreen())
	if err := program.Start(); err != nil {
		debug.Errorf("basic-space", "Application failed to start: %v", err)
		fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
		os.Exit(1)
	}

	debug.Info("basic-space", "Basic space screen test shutting down")
}
