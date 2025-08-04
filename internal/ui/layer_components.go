package ui

import (
	"github.com/charmbracelet/lipgloss/v2"
	"harvester/pkg/rendering"
	"strings"
)

// LayerTVFrame creates a TV frame using lipgloss styling instead of manual glyph generation
type LayerTVFrame struct {
	width, height int
	style         lipgloss.Style
}

func NewLayerTVFrame(w, h int) *LayerTVFrame {
	// Create a dark frame style
	frameStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#181820")).
		Background(lipgloss.Color("#181820")).
		Width(w - 6). // Account for border
		Height(h - 6)

	return &LayerTVFrame{
		width:  w,
		height: h,
		style:  frameStyle,
	}
}

func (ltf *LayerTVFrame) GetLayer() rendering.Layer { return rendering.LayerTVFrame }
func (ltf *LayerTVFrame) GetZ() int                 { return rendering.ZFrame }

func (ltf *LayerTVFrame) ToLipglossLayer() *lipgloss.Layer {
	// Create empty content with the frame style
	content := strings.Repeat(" ", ltf.width-6)
	for i := 0; i < ltf.height-6; i++ {
		if i > 0 {
			content += "\n" + strings.Repeat(" ", ltf.width-6)
		}
	}

	styledContent := ltf.style.Render(content)
	return lipgloss.NewLayer(styledContent).
		X(0).
		Y(0).
		Z(ltf.GetZ()).
		ID("tv-frame")
}


// LayerTextPanel creates a text panel with automatic styling and layout
type LayerTextPanel struct {
	title   string
	content string
	x, y    int
	width   int
	style   lipgloss.Style
	layer   rendering.Layer
	z       int
}

func NewLayerTextPanel(title, content string, width int) *LayerTextPanel {
	// Create a styled panel
	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3a3a4a")).
		Background(lipgloss.Color("#1a1a2e")).
		Foreground(lipgloss.Color("#eee8d5")).
		Padding(1, 2).
		Width(width)

	return &LayerTextPanel{
		title:   title,
		content: content,
		width:   width,
		style:   panelStyle,
		layer:   rendering.LayerUI,
		z:       rendering.ZUI,
	}
}

func (ltp *LayerTextPanel) GetLayer() rendering.Layer { return ltp.layer }
func (ltp *LayerTextPanel) GetZ() int                 { return ltp.z }

func (ltp *LayerTextPanel) WithPosition(x, y int) *LayerTextPanel {
	ltp.x, ltp.y = x, y
	return ltp
}

func (ltp *LayerTextPanel) WithLayer(layer rendering.Layer) *LayerTextPanel {
	ltp.layer = layer
	return ltp
}

func (ltp *LayerTextPanel) WithZ(z int) *LayerTextPanel {
	ltp.z = z
	return ltp
}

func (ltp *LayerTextPanel) ToLipglossLayer() *lipgloss.Layer {
	// Combine title and content
	fullContent := ltp.content
	if ltp.title != "" {
		titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#fdf6e3"))
		styledTitle := titleStyle.Render(ltp.title)
		fullContent = styledTitle + "\n\n" + ltp.content
	}

	styledContent := ltp.style.Render(fullContent)
	return lipgloss.NewLayer(styledContent).
		X(ltp.x).
		Y(ltp.y).
		Z(ltp.GetZ()).
		ID("text-panel")
}

// LayerStatusBar creates a status bar with automatic layout
type LayerStatusBar struct {
	items []StatusItem
	width int
	x, y  int
	style lipgloss.Style
	layer rendering.Layer
	z     int
}

type StatusItem struct {
	Label string
	Value string
	Style lipgloss.Style
}

func NewLayerStatusBar(width int) *LayerStatusBar {
	// Create status bar style
	barStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#16213e")).
		Foreground(lipgloss.Color("#0f3460")).
		Width(width).
		Padding(0, 1)

	return &LayerStatusBar{
		width: width,
		style: barStyle,
		layer: rendering.LayerUI,
		z:     rendering.ZUI + 1, // Above other UI elements
	}
}

func (lsb *LayerStatusBar) GetLayer() rendering.Layer { return lsb.layer }
func (lsb *LayerStatusBar) GetZ() int                 { return lsb.z }

func (lsb *LayerStatusBar) AddItem(label, value string) *LayerStatusBar {
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a8a8a8"))

	lsb.items = append(lsb.items, StatusItem{
		Label: label,
		Value: value,
		Style: labelStyle,
	})
	return lsb
}

func (lsb *LayerStatusBar) WithPosition(x, y int) *LayerStatusBar {
	lsb.x, lsb.y = x, y
	return lsb
}

func (lsb *LayerStatusBar) ToLipglossLayer() *lipgloss.Layer {
	// Build status content
	var parts []string
	for _, item := range lsb.items {
		labelStyled := item.Style.Render(item.Label + ":")
		valueStyled := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Bold(true).Render(item.Value)
		parts = append(parts, labelStyled+" "+valueStyled)
	}

	content := strings.Join(parts, " | ")
	styledContent := lsb.style.Render(content)

	return lipgloss.NewLayer(styledContent).
		X(lsb.x).
		Y(lsb.y).
		Z(lsb.GetZ()).
		ID("status-bar")
}

// LayerGameContent creates game content using lipgloss styling
type LayerGameContent struct {
	entities []GameEntity
	width    int
	height   int
	style    lipgloss.Style
}

type GameEntity struct {
	X, Y  int
	Char  string
	Style lipgloss.Style
}

func NewLayerGameContent(w, h int) *LayerGameContent {
	gameStyle := lipgloss.NewStyle().
		Width(w).
		Height(h)

	return &LayerGameContent{
		width:  w,
		height: h,
		style:  gameStyle,
	}
}

func (lgc *LayerGameContent) GetLayer() rendering.Layer { return rendering.LayerGame }
func (lgc *LayerGameContent) GetZ() int                 { return rendering.ZGame }

func (lgc *LayerGameContent) AddEntity(x, y int, char string, style lipgloss.Style) *LayerGameContent {
	lgc.entities = append(lgc.entities, GameEntity{
		X: x, Y: y, Char: char, Style: style,
	})
	return lgc
}

func (lgc *LayerGameContent) ToLipglossLayer() *lipgloss.Layer {
	// Create a 2D grid for positioning entities
	grid := make([][]string, lgc.height)
	for i := range grid {
		grid[i] = make([]string, lgc.width)
		for j := range grid[i] {
			grid[i][j] = " " // Empty space
		}
	}

	// Place entities in the grid
	for _, entity := range lgc.entities {
		if entity.X >= 0 && entity.X < lgc.width && entity.Y >= 0 && entity.Y < lgc.height {
			styledChar := entity.Style.Render(entity.Char)
			grid[entity.Y][entity.X] = styledChar
		}
	}

	// Convert grid to string
	var content strings.Builder
	for y, row := range grid {
		if y > 0 {
			content.WriteString("\n")
		}
		for _, cell := range row {
			content.WriteString(cell)
		}
	}

	return lipgloss.NewLayer(content.String()).
		X(0).
		Y(0).
		Z(lgc.GetZ()).
		ID("game-content")
}
