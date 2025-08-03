package ui

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"harvester/pkg/ecs"
)

// StatusSection represents a section of the status bar
type StatusSection struct {
	Label string
	Value string
	Style func(string) string
}

// StatusBar component with multiple sections
func StatusBarComponent(w int, sections []StatusSection) string {
	var parts []string
	for _, section := range sections {
		if section.Style != nil {
			parts = append(parts, section.Label+section.Style(section.Value))
		} else {
			parts = append(parts, section.Label+section.Value)
		}
	}
	
	content := strings.Join(parts, Muted("  |  "))
	return StatusBar(w, content)
}

// QuestPanel component
type QuestPanelData struct {
	Status string
}

func QuestPanel(data QuestPanelData) string {
	return lipgloss.JoinVertical(lipgloss.Left,
		Header("╭─ QUEST ─╮"),
		Muted("Status: ")+Highlight(data.Status),
		"",
	)
}

// ControlsPanel component
type ControlItem struct {
	Key         string
	Description string
}

type ControlsGroup struct {
	Title string
	Items []ControlItem
}

func ControlsPanel(groups []ControlsGroup) string {
	sections := []string{Header("╭─ CONTROLS ─╮")}
	
	for i, group := range groups {
		if i > 0 {
			sections = append(sections, "")
		}
		
		sections = append(sections, Muted(group.Title+":"))
		for _, item := range group.Items {
			sections = append(sections, "  "+Highlight(item.Key)+" "+Muted(item.Description))
		}
	}
	
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// MapPanel component with optional border
type MapPanelOptions struct {
	Width    int
	Height   int
	Border   bool
	Title    string
}

func MapPanel(content string, opts MapPanelOptions) string {
	panel := Sized(opts.Width, opts.Height, content)
	
	if opts.Border {
		if opts.Title != "" {
			title := Header(" " + opts.Title + " ")
			panel = lipgloss.JoinVertical(lipgloss.Center, title, Bordered(panel))
		} else {
			panel = Bordered(panel)
		}
	}
	
	return panel
}

// LogPanel component with message type styling
type LogMessage struct {
	Text string
	Type LogMessageType
}

type LogMessageType int

const (
	LogInfo LogMessageType = iota
	LogWarning
	LogError
	LogSuccess
)

func LogPanel(messages []LogMessage, width int) string {
	if len(messages) == 0 {
		return Sized(width, 0, Muted("No messages"))
	}
	
	var styledMessages []string
	for _, msg := range messages {
		var styled string
		switch msg.Type {
		case LogInfo:
			styled = Muted(msg.Text)
		case LogWarning:
			styled = Stat("", msg.Text, StatWarning)
		case LogError:
			styled = Stat("", msg.Text, StatDanger)
		case LogSuccess:
			styled = Stat("", msg.Text, StatGood)
		default:
			styled = msg.Text
		}
		styledMessages = append(styledMessages, styled)
	}
	
	content := strings.Join(styledMessages, "\n")
	return Sized(width, 0, content)
}

// PlayerStatsComponent for detailed player information
type PlayerStatsData struct {
	Fuel  int
	Hull  int
	Drive int
}

func PlayerStatsComponent(stats PlayerStatsData) string {
	return Stat("Fuel ", itoa(stats.Fuel), GetStatColor(stats.Fuel, 100)) +
		Muted("  ") +
		Stat("Hull ", itoa(stats.Hull), GetStatColor(stats.Hull, 100)) +
		Muted("  ") +
		Stat("Drive ", itoa(stats.Drive), StatGood)
}

// LocationComponent for current location display
type LocationData struct {
	Layer    string
	Planet   int
	Depth    int
}

func LocationComponent(location LocationData) string {
	return Highlight("Layer ") + Muted(location.Layer) +
		Muted("  |  ") +
		Highlight("Planet ") + Muted(itoa(location.Planet)) +
		Muted("  |  ") +
		Highlight("Depth ") + Muted(itoa(location.Depth))
}

// GameInfoComponent for tick and other game state
type GameInfoData struct {
	Tick int
}

func GameInfoComponent(info GameInfoData) string {
	return Muted("Tick " + itoa(info.Tick))
}

// Enhanced StatusBar using components and advanced styling
func EnhancedStatusBar(w int, location LocationData, stats PlayerStatsData, info GameInfoData) string {
	builder := NewComponentBuilder()
	
	left := builder.
		Add(LocationComponent(location)).
		Add(Muted("  |  ")).
		Add(GameInfoComponent(info)).
		Layout(lipgloss.Center).
		Build()
	
	right := PlayerStatsComponent(stats)
	
	// Use StyleBuilder for the final layout
	return NewStyleBuilder().
		Width(w).
		Background(GetCurrentTheme().Bg).
		Foreground(GetCurrentTheme().Text).
		PaddingHorizontal(1).
		Render(lipgloss.JoinHorizontal(lipgloss.Center,
			left,
			strings.Repeat(" ", max(0, w-lipgloss.Width(left)-lipgloss.Width(right)-4)),
			right,
		))
}

// Enhanced components using advanced styling patterns

// DynamicQuestPanel with state-aware styling
func DynamicQuestPanel(data QuestPanelData, state GameState) string {
	return DynamicPanel("QUEST", 
		"Status: "+data.Status, 
		state, 
		PanelOptions{Border: true})
}

// AnimatedStatusComponent with optional effects
type AnimatedStatusData struct {
	Label     string
	Value     string
	Status    StatStatus
	Trend     TrendDirection
	Animated  bool
	Frame     int
}

func AnimatedStatusComponent(data AnimatedStatusData) string {
	if data.Animated && data.Status == StatDanger {
		// Blinking effect for critical stats
		return BlinkingText(StatWithTrend(data.Label, data.Value, data.Trend, data.Status), data.Frame)
	}
	return StatWithTrend(data.Label, data.Value, data.Trend, data.Status)
}

// AdvancedPlayerStatsComponent with trends and animations
func AdvancedPlayerStatsComponent(stats PlayerStatsData, prevStats PlayerStatsData, frame int) string {
	builder := NewComponentBuilder()
	
	// Fuel with trend
	fuelTrend := TrendFlat
	if stats.Fuel > prevStats.Fuel {
		fuelTrend = TrendUp
	} else if stats.Fuel < prevStats.Fuel {
		fuelTrend = TrendDown
	}
	
	fuelData := AnimatedStatusData{
		Label:    "Fuel ",
		Value:    itoa(stats.Fuel),
		Status:   GetStatColor(stats.Fuel, 100),
		Trend:    fuelTrend,
		Animated: stats.Fuel < 20, // Animate when critical
		Frame:    frame,
	}
	
	// Hull with trend
	hullTrend := TrendFlat
	if stats.Hull > prevStats.Hull {
		hullTrend = TrendUp
	} else if stats.Hull < prevStats.Hull {
		hullTrend = TrendDown
	}
	
	hullData := AnimatedStatusData{
		Label:    "Hull ",
		Value:    itoa(stats.Hull),
		Status:   GetStatColor(stats.Hull, 100),
		Trend:    hullTrend,
		Animated: stats.Hull < 20, // Animate when critical
		Frame:    frame,
	}
	
	return builder.
		Add(AnimatedStatusComponent(fuelData)).
		Add(Muted("  ")).
		Add(AnimatedStatusComponent(hullData)).
		Add(Muted("  ")).
		Add(Stat("Drive ", itoa(stats.Drive), StatGood)).
		Layout(lipgloss.Left).
		Build()
}

// ThemedPanel applies theme-aware styling based on game layer
func ThemedPanel(title, content string, layer ecs.GameLayer, width, height int) string {
	themeColor := GetStateColor(layer)
	
	panel := NewComponentBuilder().
		Header(title).
		Content(content, nil).
		Layout(lipgloss.Top).
		Build()
	
	return NewStyleBuilder().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderColor(GetCurrentTheme().Border).
		Theme(themeColor).
		PaddingHorizontal(1).
		Render(panel)
}

// ResponsiveControlsPanel adapts to available space
func ResponsiveControlsPanel(groups []ControlsGroup, availableWidth int) string {
	builder := NewComponentBuilder()
	builder.Header("╭─ CONTROLS ─╮")
	
	for i, group := range groups {
		if i > 0 {
			builder.Add("")
		}
		
		builder.Content(group.Title+":", Muted)
		
		if availableWidth < 25 {
			// Compact mode - show only essential controls
			for _, item := range group.Items {
				if isEssentialControl(item.Key) {
					builder.Add("  " + Highlight(item.Key) + " " + Muted(item.Description))
				}
			}
		} else {
			// Full mode - show all controls
			for _, item := range group.Items {
				builder.Add("  " + Highlight(item.Key) + " " + Muted(item.Description))
			}
		}
	}
	
	return builder.Layout(lipgloss.Top).Build()
}

func isEssentialControl(key string) bool {
	essential := []string{"h j k l", "> ", "q "}
	for _, e := range essential {
		if key == e {
			return true
		}
	}
	return false
}

// Utility functions
func itoa(i int) string {
	return strconv.Itoa(i)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}