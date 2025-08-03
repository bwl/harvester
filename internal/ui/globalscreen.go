package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/rendering"
	"harvester/pkg/timing"
)


// ScreenType represents different screen states
type ScreenType int

const (
	ScreenStart ScreenType = iota
	ScreenSpace
	ScreenPlanet
	ScreenQuitting
)

// SubScreen interface that all screens must implement
type SubScreen interface {
	tea.Model
	HandleGlobalAction(action GlobalAction) (SubScreen, tea.Cmd)
}

// RenderableScreen interface for screens that can register content with ViewRenderer
type RenderableScreen interface {
	RegisterContent(renderer *rendering.ViewRenderer)
}

// GlobalAction represents actions that affect the entire application
type GlobalAction int

const (
	ActionNone GlobalAction = iota
	ActionStartShutdown
	ActionSwitchToGame
	ActionSwitchToStart
)

// GlobalScreen manages the overall application state and window-wide effects
type GlobalScreen struct {
	currentScreen ScreenType
	subScreen     SubScreen

	width  int
	height int

	// Global effects
	shutdownAnim  *timing.AnimationState
	openingAnim   *timing.AnimationState
	openingPending bool
	quitting      bool

	// Screen transition state
	transitioning bool
	nextScreen    ScreenType
	nextSubScreen SubScreen

	// MVP compositor at global level (used by sub-screens already)
	renderer *rendering.ViewRenderer
	fl *timing.FrameLimiter
}

func NewGlobalScreen() *GlobalScreen {
	startScreen := NewStartScreen()

	return &GlobalScreen{
		currentScreen: ScreenStart,
		subScreen:     startScreen,
		shutdownAnim:  nil,
		quitting:      false,
		transitioning: false,
	}
}

func (g *GlobalScreen) Init() tea.Cmd {
	// Initialize global timer
	timing.GetGlobalTimer()
	// Delay opening animation start by 500ms; keep hidden until it starts
	g.openingAnim = nil
	g.openingPending = true

	// Start ticker and schedule opening start; we throttle in Update using FrameLimiter
	g.fl = timing.NewFrameLimiter(60)
	return tea.Batch(
		g.subScreen.Init(),
		tea.Tick(time.Millisecond*1, func(t time.Time) tea.Msg { return t }),
		tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg { return GlobalAction(ActionSwitchToStart) }),
	)
}

func (g *GlobalScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Update global timer
	timing.UpdateGlobalTimer()

	// Handle shutdown animation
	if g.shutdownAnim != nil {
		g.shutdownAnim.Update()
		if g.shutdownAnim.IsFinished() {
			return g, tea.Quit
		}
	}
	// Kick off opening animation after delay message
	if act, ok := msg.(GlobalAction); ok && act == ActionSwitchToStart {
		g.openingAnim = timing.NewAnimation(30, false)
		g.openingPending = false
	}
	// Handle opening animation
	if g.openingAnim != nil {
		if g.width == 0 || g.height == 0 {
			// wait for size before animating
		} else {
			g.openingAnim.Update()
			if g.openingAnim.IsFinished() { g.openingAnim = nil }
		}
	}

	// Handle tick messages; throttle to 60 FPS
	if _, ok := msg.(time.Time); ok {
		if g.fl == nil { g.fl = timing.NewFrameLimiter(60) }
		if !g.fl.Allow() {
			return g, tea.Tick(time.Millisecond*1, func(t time.Time) tea.Msg { return t })
		}
		return g, tea.Tick(time.Second/60, func(t time.Time) tea.Msg { return t })
	}

	// Handle window resize at global level
	if windowMsg, ok := msg.(tea.WindowSizeMsg); ok {
		g.width = windowMsg.Width
		g.height = windowMsg.Height
		if gs, ok := g.subScreen.(*StartScreen); ok {
			_ = gs // StartScreen handles its own renderer sizing
		}
	}

	// Handle global quit keys
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "ctrl+c" {
			// Immediate quit on Ctrl+C
			return g, tea.Quit
		}

		if keyMsg.String() == "q" && !g.isInGameScreen() {
			// Start shutdown animation for Q (only when not in game)
			g.startShutdownAnimation()
			return g, nil
		}
	}

	// If we're animating shutdown, ignore input to sub-screens
	if g.shutdownAnim != nil {
		return g, nil
	}

	// Handle screen transitions
	if g.transitioning {
		g.completeTransition()
		return g, g.subScreen.Init()
	}

	// Forward message to current sub-screen
	newSubScreen, cmd := g.subScreen.Update(msg)

	// Handle type assertion for SubScreen interface
	if subScreen, ok := newSubScreen.(SubScreen); ok {
		g.subScreen = subScreen
	}

	// Check if sub-screen wants to trigger global actions
	if g.currentScreen == ScreenStart {
		if startScreen, ok := g.subScreen.(*StartScreen); ok {
			result := startScreen.GetResult()
			if result != nil {
				return g.handleStartScreenResult(result)
			}
		}
	}

	return g, cmd
}

func (g *GlobalScreen) View() string {
	if g.width == 0 || g.height == 0 {
		return "Loading..."
	}

	// Hidden until opening starts
	if g.openingPending {
		return joinLines(make([]string, g.height))
	}

	// Return sub-screen content - effects will be handled by RootView through RegisterContent
	return g.subScreen.View()
}

// RegisterContent implements RenderableScreen interface
func (g *GlobalScreen) RegisterContent(renderer *rendering.ViewRenderer) {
	// Hidden until opening starts
	if g.openingPending {
		return
	}

	// Let sub-screen register its content first
	if renderableScreen, ok := g.subScreen.(RenderableScreen); ok {
		renderableScreen.RegisterContent(renderer)
	} else {
		// Fallback: convert string-based screen to renderable content
		content := g.subScreen.View()
		if content != "" {
			textBlock := newTextBlock(content, g.width, g.height)
			renderer.RegisterContent(textBlock)
		}
	}

	// Add global effect overlays
	if g.shutdownAnim != nil {
		progress := g.shutdownAnim.Progress()
		shutdownOverlay := NewCRTShutdownOverlay(g.width, g.height, progress)
		renderer.RegisterContent(shutdownOverlay)
	}

	if g.openingAnim != nil {
		progress := g.openingAnim.Progress()
		openingOverlay := NewCRTOpeningOverlay(g.width, g.height, progress)
		renderer.RegisterContent(openingOverlay)
	}
}

func (g *GlobalScreen) isInGameScreen() bool {
	return g.currentScreen == ScreenSpace || g.currentScreen == ScreenPlanet
}

func (g *GlobalScreen) startShutdownAnimation() {
	// 500ms shutdown using harmonica spring
	g.shutdownAnim = timing.NewAnimation(30, false)
	g.quitting = true
}

func (g *GlobalScreen) handleStartScreenResult(result *StartResult) (tea.Model, tea.Cmd) {
	switch result.Action {
	case ActionQuit:
		g.startShutdownAnimation()
		return g, nil

	case ActionContinue, ActionLoadSlot, ActionNewGame:
		// Transition to space screen (games start in space)
		return g.transitionToSpace(result)

	default:
		return g, nil
	}
}

func (g *GlobalScreen) transitionToSpace(startResult *StartResult) (tea.Model, tea.Cmd) {
	// Create new space screen based on start result
	spaceScreen := g.createSpaceScreen(startResult)

	g.nextScreen = ScreenSpace
	g.nextSubScreen = spaceScreen
	g.transitioning = true

	return g, nil
}

func (g *GlobalScreen) transitionToPlanet(startResult *StartResult) (tea.Model, tea.Cmd) {
	// Create new planet screen based on start result
	planetScreen := g.createPlanetScreen(startResult)

	g.nextScreen = ScreenPlanet
	g.nextSubScreen = planetScreen
	g.transitioning = true

	return g, nil
}

func (g *GlobalScreen) createSpaceScreen(result *StartResult) SubScreen {
	// Create space navigation screen
	return NewSpaceScreen(result)
}

func (g *GlobalScreen) createPlanetScreen(result *StartResult) SubScreen {
	// Create planet exploration screen
	return NewPlanetScreen(result)
}

func (g *GlobalScreen) completeTransition() {
	g.currentScreen = g.nextScreen
	g.subScreen = g.nextSubScreen
	g.transitioning = false
	g.nextSubScreen = nil
}


// Helper functions for line manipulation
func splitLines(content string) []string {
	if content == "" {
		return []string{}
	}
	lines := []string{}
	current := ""
	for _, char := range content {
		if char == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func joinLines(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	result := lines[0]
	for i := 1; i < len(lines); i++ {
		result += "\n" + lines[i]
	}
	return result
}
