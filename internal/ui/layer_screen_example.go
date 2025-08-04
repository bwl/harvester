package ui

import (
	"github.com/charmbracelet/lipgloss/v2"
	"harvester/pkg/rendering"
)

// LayerStartScreen shows how to convert a screen to use Layer-based rendering
type LayerStartScreen struct {
	width, height int
	title         string
	menuItems     []string
	selectedIndex int
}

func NewLayerStartScreen(w, h int) *LayerStartScreen {
	return &LayerStartScreen{
		width:         w,
		height:        h,
		title:         "🌌 HARVESTER",
		menuItems:     []string{"New Game", "Load Game", "Settings", "Quit"},
		selectedIndex: 0,
	}
}

// RegisterContent shows the new approach for screen registration
func (lss *LayerStartScreen) RegisterContent(renderer *rendering.CanvasRenderer) {
	// Create background
	background := NewLayerBackground(lss.width, lss.height)
	renderer.RegisterContent(background)

	// Create TV frame
	frame := NewLayerTVFrame(lss.width, lss.height)
	renderer.RegisterContent(frame)

	// Create title
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffd700")).
		Background(lipgloss.Color("#1a1a2e")).
		Bold(true).
		Align(lipgloss.Center).
		Width(lss.width-20).
		Padding(1, 0).
		Margin(2, 0)

	titlePanel := rendering.NewStyledContent(
		rendering.LayerUI,
		rendering.ZUI,
		lss.title,
	).WithStyle(titleStyle).WithPosition(10, 3).WithID("title")

	renderer.RegisterContent(titlePanel)

	// Create menu
	lss.createMenu(renderer)

	// Create footer
	footerText := "Use ↑↓ to navigate, Enter to select, ESC to quit"
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Align(lipgloss.Center).
		Width(lss.width - 20)

	footer := rendering.NewStyledContent(
		rendering.LayerUI,
		rendering.ZUI,
		footerText,
	).WithStyle(footerStyle).WithPosition(10, lss.height-5).WithID("footer")

	renderer.RegisterContent(footer)
}

func (lss *LayerStartScreen) createMenu(renderer *rendering.CanvasRenderer) {
	menuY := lss.height/2 - len(lss.menuItems)/2

	for i, item := range lss.menuItems {
		var itemStyle lipgloss.Style

		if i == lss.selectedIndex {
			// Selected item style
			itemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#ffd700")).
				Bold(true).
				Padding(0, 2).
				Width(20).
				Align(lipgloss.Center)
		} else {
			// Normal item style
			itemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#cccccc")).
				Background(lipgloss.Color("#1a1a2e")).
				Padding(0, 2).
				Width(20).
				Align(lipgloss.Center)
		}

		menuItem := rendering.NewStyledContent(
			rendering.LayerMenu,
			rendering.ZMenu,
			item,
		).WithStyle(itemStyle).
			WithPosition(lss.width/2-10, menuY+i*2).
			WithID("menu-item-" + item)

		renderer.RegisterContent(menuItem)
	}
}

// SetSelectedIndex updates the selected menu item
func (lss *LayerStartScreen) SetSelectedIndex(index int) {
	if index >= 0 && index < len(lss.menuItems) {
		lss.selectedIndex = index
	}
}

// LayerGameScreen shows how to create a game screen with Layer-based rendering
type LayerGameScreen struct {
	width, height    int
	playerX, playerY int
	entities         []GameEntity
}

func NewLayerGameScreen(w, h int) *LayerGameScreen {
	return &LayerGameScreen{
		width:   w,
		height:  h,
		playerX: w / 2,
		playerY: h / 2,
	}
}

func (lgs *LayerGameScreen) RegisterContent(renderer *rendering.CanvasRenderer) {
	// Background
	background := NewLayerBackground(lgs.width, lgs.height)
	renderer.RegisterContent(background)

	// Game content area
	gameArea := rendering.NewStyledContent(
		rendering.LayerGame,
		rendering.ZContent,
		lgs.createGameMap(),
	).WithPosition(5, 3).WithID("game-map")

	renderer.RegisterContent(gameArea)

	// Player
	playerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff6b6b")).
		Bold(true)

	player := rendering.NewStyledContent(
		rendering.LayerGame,
		rendering.ZGame+10, // Above other game elements
		"@",
	).WithStyle(playerStyle).
		WithPosition(lgs.playerX, lgs.playerY).
		WithID("player")

	renderer.RegisterContent(player)

	// HUD
	lgs.createHUD(renderer)
}

func (lgs *LayerGameScreen) createGameMap() string {
	// Simple procedural map
	mapContent := ""
	for y := 0; y < 15; y++ {
		for x := 0; x < 60; x++ {
			if x == 0 || x == 59 || y == 0 || y == 14 {
				mapContent += "#" // Walls
			} else if (x+y)%7 == 0 {
				mapContent += "." // Scattered dots
			} else if (x*y)%13 == 0 {
				mapContent += "*" // Resources
			} else {
				mapContent += " " // Empty space
			}
		}
		if y < 14 {
			mapContent += "\n"
		}
	}
	return mapContent
}

func (lgs *LayerGameScreen) createHUD(renderer *rendering.CanvasRenderer) {
	// Health bar
	healthStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff6b6b")).
		Background(lipgloss.Color("#1a1a2e")).
		Padding(0, 1)

	health := rendering.NewStyledContent(
		rendering.LayerHUD,
		rendering.ZHUD,
		"Health: ████████░░ 80%",
	).WithStyle(healthStyle).
		WithPosition(5, 1).
		WithID("health")

	renderer.RegisterContent(health)

	// Score
	scoreStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4ecdc4")).
		Background(lipgloss.Color("#1a1a2e")).
		Padding(0, 1)

	score := rendering.NewStyledContent(
		rendering.LayerHUD,
		rendering.ZHUD,
		"Score: 1,337",
	).WithStyle(scoreStyle).
		WithPosition(30, 1).
		WithID("score")

	renderer.RegisterContent(score)

	// Mini-map
	minimapStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3a3a4a")).
		Background(lipgloss.Color("#0a0a0f")).
		Foreground(lipgloss.Color("#666666")).
		Width(12).
		Height(6).
		Padding(0, 1)

	minimap := rendering.NewStyledContent(
		rendering.LayerHUD,
		rendering.ZHUD,
		"╔══════╗\n║  @   ║\n║      ║\n║   *  ║\n╚══════╝",
	).WithStyle(minimapStyle).
		WithPosition(lgs.width-15, 1).
		WithID("minimap")

	renderer.RegisterContent(minimap)
}

func (lgs *LayerGameScreen) MovePlayer(dx, dy int) {
	newX := lgs.playerX + dx
	newY := lgs.playerY + dy

	// Simple bounds checking
	if newX >= 6 && newX < lgs.width-10 && newY >= 4 && newY < lgs.height-5 {
		lgs.playerX = newX
		lgs.playerY = newY
	}
}
