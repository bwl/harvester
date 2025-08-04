package ui

import (
	"fmt"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	screens "harvester/internal/ui/screens"
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"harvester/pkg/rendering"
	"harvester/pkg/systems"
	"harvester/pkg/timing"
)

type StartAction int

const (
	ActionContinue StartAction = iota
	ActionLoadSlot
	ActionNewGame
	ActionQuit
)

type StartResult struct {
	Action  StartAction
	SlotNum int // for ActionLoadSlot
}

type StartScreen struct {
	selected     int
	menuItems    []string
	saveSlots    []SaveSlotInfo
	showSlots    bool
	selectedSlot int

	width  int
	height int

	result *StartResult

	// Animation state
	shutdownAnim *timing.AnimationState

	// Background terrain using TOFT system
	backgroundWorld *ecs.World
	renderer        *systems.Render

	// Save game manager
	saveManager *SaveGameManager
}

func NewStartScreen() *StartScreen {
	s := &StartScreen{
		selected:     0,
		showSlots:    false,
		selectedSlot: 0,
		shutdownAnim: nil,
		saveManager:  NewSaveGameManager(),
	}

	s.scanSaveFiles()
	s.buildMenuItems()
	s.initBackgroundTerrain()

	return s
}

func (s *StartScreen) initBackgroundTerrain() {
	// Create a background world for terrain generation
	r := rand.New(rand.NewSource(42)) // Fixed seed for consistent background
	s.backgroundWorld = ecs.NewWorld(r)

	// Set world context for TOFT planet surface
	ctx := ecs.WorldContext{
		CurrentLayer: ecs.LayerPlanetSurface,
		PlanetID:     42, // Fixed planet ID for consistent background
		Depth:        0,
		BiomeType:    0,
	}
	ecs.SetWorldContext(s.backgroundWorld, ctx)

	// Add WorldInfo with initial dimensions (will be updated when screen size is known)
	worldInfo := components.WorldInfo{
		Width:  80, // Default width
		Height: 24, // Default height
	}
	ecs.Add(s.backgroundWorld, 1, worldInfo)

	// Initialize map renderer using the same system as the game
	s.renderer = &systems.Render{}
}

func (s *StartScreen) generateBackgroundTerrain() {
	if s.backgroundWorld == nil || s.width == 0 || s.height == 0 {
		return
	}

	// Update world size to match screen dimensions
	worldInfo := components.WorldInfo{
		Width:  s.width,
		Height: s.height,
	}
	ecs.Add(s.backgroundWorld, 1, worldInfo)

	// Clear existing terrain entities
	var toRemove []ecs.Entity
	ecs.View1Of[components.Tile](s.backgroundWorld).Each(func(e ecs.Entity, _ *components.Tile) {
		toRemove = append(toRemove, e)
	})
	for _, e := range toRemove {
		s.backgroundWorld.Destroy(e)
	}

	// Use the same TerrainGen system as TOFT
	terrainGen := systems.TerrainGen{}
	terrainGen.Update(0, s.backgroundWorld)
}

// Compositor-driven background content
func (s *StartScreen) renderBackgroundContent() rendering.RenderableContent {
	if s.backgroundWorld == nil || s.renderer == nil || s.width == 0 || s.height == 0 {
		return nil
	}
	s.renderer.Update(0, s.backgroundWorld)
	glyphs := make([][]rendering.Glyph, s.height)
	for y := 0; y < s.height; y++ {
		row := make([]rendering.Glyph, s.width)
		for x := 0; x < s.width; x++ {
			row[x] = rendering.Glyph{Char: ' '}
		}
		glyphs[y] = row
	}
	for _, drawable := range s.renderer.Output {
		x, y := drawable.X, drawable.Y
		if x >= 0 && x < s.width && y >= 0 && y < s.height {
			glyphs[y][x] = rendering.Glyph{Char: rune(drawable.Glyph)}
		}
	}
	return screens.NewTerrainContent(glyphs, s.width, s.height)
}

// Compositor-driven menu content
func (s *StartScreen) renderMenuContent() rendering.RenderableContent {
	var content string
	if s.showSlots {
		content = s.renderSlotSelection()
	} else {
		content = s.renderMainMenu()
	}
	lines := strings.Split(content, "\n")
	glyphs := rendering.RenderLipglossString(lines, rendering.Color{R: 220, G: 220, B: 220}, rendering.Color{}, rendering.StyleNone)
	w := 0
	for _, l := range lines {
		if lw := len([]rune(l)); lw > w {
			w = lw
		}
	}
	h := len(lines)
	return screens.NewMenuContent(glyphs, w, h)
}

// Simple ANSI code stripper to measure actual text width
func stripANSI(s string) string {
	result := ""
	inEscape := false
	for _, char := range s {
		if char == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if char == 'm' {
				inEscape = false
			}
			continue
		}
		result += string(char)
	}
	return result
}

func (s *StartScreen) scanSaveFiles() {
	// Check autosave
	autosaveExists := s.saveManager.HasAutosave()

	// Check save slots
	s.saveSlots = s.saveManager.GetSaveSlots()

	// Build menu based on what saves exist
	s.menuItems = []string{}
	if autosaveExists {
		s.menuItems = append(s.menuItems, "Continue")
	}

	// Only show "Load Game" if any slots exist
	hasSlots := false
	for _, slot := range s.saveSlots {
		if slot.Exists {
			hasSlots = true
			break
		}
	}
	if hasSlots {
		s.menuItems = append(s.menuItems, "Load Game")
	}

	s.menuItems = append(s.menuItems, "New Game", "Quit")
}

func (s *StartScreen) buildMenuItems() {
	// Menu items are built in scanSaveFiles
}

func (s *StartScreen) Init() tea.Cmd {
	// Initialize global timer if needed
	timing.GetGlobalTimer()
	return nil
}

func (s *StartScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Update global timer
	timing.UpdateGlobalTimer()

	// Update shutdown animation if running
	if s.shutdownAnim != nil {
		s.shutdownAnim.Update()
		if s.shutdownAnim.IsFinished() {
			// Animation complete, quit
			return s, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.generateBackgroundTerrain() // Regenerate terrain for new size

	case tea.KeyMsg:
		if s.shutdownAnim != nil {
			// Ignore input during animation
			return s, nil
		}

		if s.showSlots {
			return s.updateSlotSelection(msg)
		}
		return s.updateMainMenu(msg)
	}

	return s, nil
}

func (s *StartScreen) HandleInput(a InputAction) tea.Cmd {
	if s.showSlots {
		switch a.Kind {
		case InputMenuUp:
			for i := s.selectedSlot - 1; i >= 0; i-- {
				if s.saveSlots[i].Exists {
					s.selectedSlot = i
					break
				}
			}
		case InputMenuDown:
			for i := s.selectedSlot + 1; i < len(s.saveSlots); i++ {
				if s.saveSlots[i].Exists {
					s.selectedSlot = i
					break
				}
			}
		case InputMenuSelect:
			if s.saveSlots[s.selectedSlot].Exists {
				s.result = &StartResult{Action: ActionLoadSlot, SlotNum: s.saveSlots[s.selectedSlot].SlotNum}
			}
		case InputMenuBack:
			s.showSlots = false
		case InputQuit:
			s.result = &StartResult{Action: ActionQuit}
		}
		return nil
	}
	switch a.Kind {
	case InputMenuUp:
		if s.selected > 0 {
			s.selected--
		}
	case InputMenuDown:
		if s.selected < len(s.menuItems)-1 {
			s.selected++
		}
	case InputMenuSelect:
		selectedItem := s.menuItems[s.selected]
		switch selectedItem {
		case "Continue":
			s.result = &StartResult{Action: ActionContinue}
		case "Load Game":
			s.showSlots = true
			s.selectedSlot = 0
			for i, slot := range s.saveSlots {
				if slot.Exists {
					s.selectedSlot = i
					break
				}
			}
		case "New Game":
			s.result = &StartResult{Action: ActionNewGame}
		case "Quit":
			s.result = &StartResult{Action: ActionQuit}
		}
	case InputMenuBack, InputQuit:
		s.result = &StartResult{Action: ActionQuit}
	}
	return nil
}

func (s *StartScreen) updateMainMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if s.selected > 0 {
			s.selected--
		}
	case "down", "j":
		if s.selected < len(s.menuItems)-1 {
			s.selected++
		}
	case "enter", " ":
		selectedItem := s.menuItems[s.selected]
		switch selectedItem {
		case "Continue":
			s.result = &StartResult{Action: ActionContinue}
			return s, tea.Quit
		case "Load Game":
			s.showSlots = true
			s.selectedSlot = 0
			// Find first existing slot to select
			for i, slot := range s.saveSlots {
				if slot.Exists {
					s.selectedSlot = i
					break
				}
			}
		case "New Game":
			s.result = &StartResult{Action: ActionNewGame}
			return s, tea.Quit
		case "Quit":
			s.result = &StartResult{Action: ActionQuit}
			// Don't start animation here - let GlobalScreen handle it
			return s, tea.Quit
		}
	case "q", "esc":
		s.result = &StartResult{Action: ActionQuit}
		// Don't start animation here - let GlobalScreen handle it
		return s, tea.Quit
	}

	return s, nil
}

func (s *StartScreen) updateSlotSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		// Find previous existing slot
		for i := s.selectedSlot - 1; i >= 0; i-- {
			if s.saveSlots[i].Exists {
				s.selectedSlot = i
				break
			}
		}
	case "down", "j":
		// Find next existing slot
		for i := s.selectedSlot + 1; i < len(s.saveSlots); i++ {
			if s.saveSlots[i].Exists {
				s.selectedSlot = i
				break
			}
		}
	case "enter", " ":
		if s.saveSlots[s.selectedSlot].Exists {
			s.result = &StartResult{
				Action:  ActionLoadSlot,
				SlotNum: s.saveSlots[s.selectedSlot].SlotNum,
			}
			return s, tea.Quit
		}
	case "esc", "backspace":
		s.showSlots = false
	case "q":
		s.result = &StartResult{Action: ActionQuit}
		// Don't start animation here - let GlobalScreen handle it
		return s, tea.Quit
	}

	return s, nil
}

func (s *StartScreen) View() string {
	return ""
}

func (s *StartScreen) renderMainMenu() string {
	// Create compact menu content without full-width styling
	title := NewStyleBuilder().
		Theme(ThemePrimary).
		Bold(true).
		Render("BUBBLE ROUGE")

	subtitle := NewStyleBuilder().
		Theme(ThemeMuted).
		Render("A Big Bang Roguelike")

	// Menu items - compact without full width
	var menuLines []string
	for i, item := range s.menuItems {
		var menuContent string
		var style *StyleBuilder

		if i == s.selected {
			// Selected item - highlighted with arrow
			menuContent = fmt.Sprintf("▶ %s ◀", item)
			style = NewStyleBuilder().Theme(ThemeSecondary).Bold(true)
		} else {
			// Unselected item
			menuContent = fmt.Sprintf("  %s  ", item)
			style = NewStyleBuilder().Theme(ThemeMuted)
		}

		menuLines = append(menuLines, style.Render(menuContent))
	}

	// Controls help - compact
	controls := NewStyleBuilder().
		Theme(ThemeMuted).
		Render("↑↓ navigate • Enter select • Q quit")

	// Layout everything vertically with spacing - no full width
	content := NewComponentBuilder().
		Add(title).
		Add("").
		Add(subtitle).
		Add("").
		Add("").
		Add(strings.Join(menuLines, "\n")).
		Add("").
		Add("").
		Add(controls).
		Layout(lipgloss.Top).
		Build()

	return content
}

func (s *StartScreen) renderSlotSelection() string {
	// Create compact slot selection content
	title := NewStyleBuilder().
		Theme(ThemePrimary).
		Bold(true).
		Render("LOAD GAME")

	// Save slot list - compact without full width
	var slotLines []string
	for i, slot := range s.saveSlots {
		if !slot.Exists {
			continue // Skip non-existent slots
		}

		timeStr := slot.ModTime.Format("2006-01-02 15:04")
		sizeStr := formatFileSize(slot.Size)

		slotContent := fmt.Sprintf("Slot %d - %s (%s) - %s",
			slot.SlotNum, slot.GameInfo, sizeStr, timeStr)

		var style *StyleBuilder
		if i == s.selectedSlot {
			// Selected slot
			slotContent = fmt.Sprintf("▶ %s ◀", slotContent)
			style = NewStyleBuilder().Theme(ThemeSecondary).Bold(true)
		} else {
			// Unselected slot
			slotContent = fmt.Sprintf("  %s  ", slotContent)
			style = NewStyleBuilder().Theme(ThemeMuted)
		}

		slotLines = append(slotLines, style.Render(slotContent))
	}

	// Controls help - compact
	controls := NewStyleBuilder().
		Theme(ThemeMuted).
		Render("↑↓ navigate • Enter load • Esc back • Q quit")

	// Layout everything - no full width
	content := NewComponentBuilder().
		Add(title).
		Add("").
		Add(strings.Join(slotLines, "\n")).
		Add("").
		Add(controls).
		Layout(lipgloss.Top).
		Build()

	return content
}

func formatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(size)/1024)
	} else {
		return fmt.Sprintf("%.1fMB", float64(size)/(1024*1024))
	}
}

func (s *StartScreen) startShutdownAnimation() {
	// Create a 30-frame shutdown animation (about 1.5 seconds at 20fps)
	s.shutdownAnim = timing.NewAnimation(30, false)
}

func (s *StartScreen) renderShutdownEffect(content string) string {
	if s.shutdownAnim == nil {
		return content
	}

	// Get animation progress (0.0 to 1.0)
	progress := s.shutdownAnim.Progress()

	// Apply easing for more natural CRT effect
	easedProgress := timing.EaseIn(progress)

	// Calculate how much of the screen to show
	// At 0.0: full height, at 1.0: collapsed to center line
	originalHeight := s.height
	collapsedHeight := int(float64(originalHeight) * (1.0 - easedProgress))

	// Minimum height of 1 to show the center line
	if collapsedHeight < 1 {
		collapsedHeight = 1
	}

	// Calculate margins to center the collapsed content
	topMargin := (originalHeight - collapsedHeight) / 2
	bottomMargin := originalHeight - collapsedHeight - topMargin

	// Split content into lines
	lines := strings.Split(content, "\n")

	// Calculate which lines to show (center portion)
	totalLines := len(lines)
	if totalLines == 0 {
		return ""
	}

	// Show middle portion of content
	startLine := int(float64(totalLines) * easedProgress / 2)
	endLine := totalLines - startLine

	if startLine >= endLine {
		// Show just the center line
		centerLine := totalLines / 2
		if centerLine < totalLines {
			lines = []string{lines[centerLine]}
		} else {
			lines = []string{""}
		}
	} else {
		lines = lines[startLine:endLine]
	}

	// Pad with empty lines to maintain terminal positioning
	var result []string

	// Add top margin
	for i := 0; i < topMargin; i++ {
		result = append(result, "")
	}

	// Add visible content
	result = append(result, lines...)

	// Add bottom margin
	for i := 0; i < bottomMargin; i++ {
		result = append(result, "")
	}

	return strings.Join(result, "\n")
}

func (s *StartScreen) HandleGlobalAction(action GlobalAction) (SubScreen, tea.Cmd) {
	switch action {
	case ActionStartShutdown:
		// Start screen can handle shutdown
		return s, nil
	default:
		return s, nil
	}
}

func (s *StartScreen) GetResult() *StartResult {
	return s.result
}

// RegisterContent implements RenderableScreen interface
func (s *StartScreen) RegisterContent(renderer *rendering.ViewRenderer) {
	// Register background terrain content
	if bg := s.renderBackgroundContent(); bg != nil {
		renderer.RegisterContent(bg)
	}

	// Register menu content
	if menu := s.renderMenuContent(); menu != nil {
		renderer.RegisterContent(menu)
	}
}
