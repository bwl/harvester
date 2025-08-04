package ui

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"harvester/pkg/ecs"
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

// ResizableScreen interface for screens that need window dimension updates
type ResizableScreen interface {
	SetDimensions(width, height int)
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
	shutdownAnim   *timing.AnimationState
	openingAnim    *timing.AnimationState
	openingPending bool
	quitting       bool

	// Screen transition state
	transitioning bool
	nextScreen    ScreenType
	nextSubScreen SubScreen

	// Frame limiter only; renderer owned by RootView
	fl *timing.FrameLimiter

	// Save game manager
	saveManager *SaveGameManager
}

func NewGlobalScreen() *GlobalScreen {
	startScreen := NewStartScreen()

	return &GlobalScreen{
		currentScreen: ScreenStart,
		subScreen:     startScreen,
		shutdownAnim:  nil,
		quitting:      false,
		transitioning: false,
		saveManager:   NewSaveGameManager(),
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
			if g.openingAnim.IsFinished() {
				g.openingAnim = nil
			}
		}
	}

	// Handle tick messages; throttle to 60 FPS
	if _, ok := msg.(time.Time); ok {
		if g.fl == nil {
			g.fl = timing.NewFrameLimiter(60)
		}
		if !g.fl.Allow() {
			return g, tea.Tick(time.Millisecond*1, func(t time.Time) tea.Msg { return t })
		}
		return g, tea.Tick(time.Second/60, func(t time.Time) tea.Msg { return t })
	}

	// Handle window resize at global level
	if windowMsg, ok := msg.(tea.WindowSizeMsg); ok {
		g.width = windowMsg.Width
		g.height = windowMsg.Height
		// Forward window size to all resizable screens
		if resizable, ok := g.subScreen.(ResizableScreen); ok {
			resizable.SetDimensions(windowMsg.Width, windowMsg.Height)
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

	// Screen transitions are now handled immediately in transition methods

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
	if g.currentScreen == ScreenSpace {
		if space, ok := g.subScreen.(*SpaceScreen); ok {
			ctx := ecs.GetWorldContext(space.model.World())
			if ctx.CurrentLayer == ecs.LayerPlanetSurface {
				return g.transitionToPlanet(&StartResult{Action: ActionContinue})
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

	// Let sub-screen register its content first (map layers/panels)
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

func (g *GlobalScreen) SetDimensions(width, height int) {
	g.width = width
	g.height = height
	
	// Forward dimensions to current sub-screen if it supports resizing
	if resizable, ok := g.subScreen.(ResizableScreen); ok {
		resizable.SetDimensions(width, height)
	}
}

func (g *GlobalScreen) HandleInput(a InputAction) tea.Cmd {
	if a.Kind == InputQuit {
		g.startShutdownAnimation()
		return nil
	}
	if ih, ok := any(g.subScreen).(InputHandler); ok {
		return ih.HandleInput(a)
	}
	return nil
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

	// Complete transition immediately to avoid requiring extra keypress
	g.completeTransition()
	return g, g.subScreen.Init()
}

func (g *GlobalScreen) transitionToPlanet(startResult *StartResult) (tea.Model, tea.Cmd) {
	// Create new planet screen based on start result
	planetScreen := g.createPlanetScreen(startResult)

	g.nextScreen = ScreenPlanet
	g.nextSubScreen = planetScreen
	g.transitioning = true

	// Complete transition immediately to avoid requiring extra keypress
	g.completeTransition()
	return g, g.subScreen.Init()
}

func (g *GlobalScreen) createSpaceScreen(result *StartResult) SubScreen {
	// Create model with appropriate save data loaded
	model := g.createModelWithSaveData(result)
	// Create space navigation screen
	spaceScreen := NewSpaceScreen(model)

	// Forward current window dimensions to the new screen
	if g.width > 0 && g.height > 0 {
		spaceScreen.SetDimensions(g.width, g.height)
	}

	return spaceScreen
}

func (g *GlobalScreen) createPlanetScreen(result *StartResult) SubScreen {
	// Create model with appropriate save data loaded
	model := g.createModelWithSaveData(result)
	// Create planet exploration screen
	planetScreen := NewPlanetScreen(model)

	// Forward current window dimensions to the new screen
	if g.width > 0 && g.height > 0 {
		planetScreen.SetDimensions(g.width, g.height)
	}

	return planetScreen
}

func (g *GlobalScreen) createModelWithSaveData(result *StartResult) *Model {
	// Create the game model with random seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	model := NewModelWithRNG(r)

	// Load save data based on start result
	switch result.Action {
	case ActionContinue:
		// Load autosave
		err := g.saveManager.LoadAutosave(model.World())
		if err != nil {
			// Log error but continue with new game
		}

	case ActionLoadSlot:
		// Load specific slot
		err := g.saveManager.LoadSlot(result.SlotNum, model.World())
		if err != nil {
			// Log error but continue with new game
		}

	case ActionNewGame:
		// Start fresh - no loading needed, starts in LayerSpace by default
	}

	return &model
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
