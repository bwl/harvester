package ui

import (
	"github.com/charmbracelet/lipgloss/v2"
	"image/color"
)

type StyleTheme struct {
	Primary       color.Color
	Secondary     color.Color
	Accent        color.Color
	Muted         color.Color
	Bg            color.Color
	Text          color.Color
	Border        color.Color
	Success       color.Color
	Warning       color.Color
	Error         color.Color
	TextSecondary color.Color
}

var theme = StyleTheme{
	Primary:       lipgloss.Color("45"),
	Secondary:     lipgloss.Color("220"),
	Accent:        lipgloss.Color("213"),
	Muted:         lipgloss.Color("240"),
	Bg:            lipgloss.Color("0"),
	Text:          lipgloss.Color("15"),
	Border:        lipgloss.Color("240"),
	Success:       lipgloss.Color("42"),
	Warning:       lipgloss.Color("214"),
	Error:         lipgloss.Color("196"),
	TextSecondary: lipgloss.Color("244"),
}

var lightTheme = StyleTheme{
	Primary:       lipgloss.Color("27"),
	Secondary:     lipgloss.Color("60"),
	Accent:        lipgloss.Color("57"),
	Muted:         lipgloss.Color("244"),
	Bg:            lipgloss.Color("255"),
	Text:          lipgloss.Color("0"),
	Border:        lipgloss.Color("244"),
	Success:       lipgloss.Color("34"),
	Warning:       lipgloss.Color("214"),
	Error:         lipgloss.Color("160"),
	TextSecondary: lipgloss.Color("240"),
}

type stylesDef struct {
	Space     func(string) string
	Header    lipgloss.Style
	Panel     lipgloss.Style
	Bordered  lipgloss.Style
	Muted     lipgloss.Style
	Highlight lipgloss.Style
}

var styles = stylesDef{}

func init() { rebuildStyles() }

func rebuildStyles() {
	styles = stylesDef{
		Space:     func(s string) string { return lipgloss.NewStyle().Foreground(theme.Muted).Render(s) },
		Header:    lipgloss.NewStyle().Bold(true).Foreground(theme.Primary),
		Panel:     lipgloss.NewStyle().Padding(0, 1).Background(theme.Bg).Foreground(theme.Text),
		Bordered:  lipgloss.NewStyle().BorderForeground(theme.Border).Border(lipgloss.NormalBorder()).Padding(0, 1),
		Muted:     lipgloss.NewStyle().Foreground(theme.Muted),
		Highlight: lipgloss.NewStyle().Foreground(theme.Secondary).Bold(true),
	}
}

func UseDarkTheme() {
	theme = StyleTheme{
		Primary:       lipgloss.Color("45"),
		Secondary:     lipgloss.Color("220"),
		Accent:        lipgloss.Color("213"),
		Muted:         lipgloss.Color("240"),
		Bg:            lipgloss.Color("0"),
		Text:          lipgloss.Color("15"),
		Border:        lipgloss.Color("240"),
		Success:       lipgloss.Color("42"),
		Warning:       lipgloss.Color("214"),
		Error:         lipgloss.Color("196"),
		TextSecondary: lipgloss.Color("244"),
	}
	rebuildStyles()
}

func UseLightTheme() {
	theme = lightTheme
	rebuildStyles()
}

func GetCurrentTheme() StyleTheme {
	return theme
}

func SetCustomTheme(customTheme StyleTheme) {
	theme = customTheme
	rebuildStyles()
}

// Style helper functions
func Panel(content string) string    { return styles.Panel.Render(content) }
func Bordered(content string) string { return styles.Bordered.Render(content) }
func Sized(w, h int, content string) string {
	return lipgloss.NewStyle().Width(w).Height(h).Render(content)
}
func Header(text string) string    { return styles.Header.Render(text) }
func Muted(text string) string     { return styles.Muted.Render(text) }
func Highlight(text string) string { return styles.Highlight.Render(text) }

// Stat status for color-coded stats
type StatStatus int

const (
	StatGood StatStatus = iota
	StatWarning
	StatDanger
)

func Stat(label, value string, status StatStatus) string {
	labelStyle := styles.Highlight
	var valueStyle lipgloss.Style

	switch status {
	case StatGood:
		valueStyle = lipgloss.NewStyle().Foreground(theme.Success)
	case StatWarning:
		valueStyle = lipgloss.NewStyle().Foreground(theme.Warning)
	case StatDanger:
		valueStyle = lipgloss.NewStyle().Foreground(theme.Error)
	default:
		valueStyle = styles.Muted
	}

	return labelStyle.Render(label) + valueStyle.Render(value)
}

// Get stat color based on percentage
func GetStatColor(current, max int) StatStatus {
	percentage := float64(current) / float64(max)
	if percentage > 0.7 {
		return StatGood
	} else if percentage > 0.3 {
		return StatWarning
	} else {
		return StatDanger
	}
}

// Advanced helpers
func PanelWithBorder(content string) string {
	return styles.Bordered.Render(styles.Panel.Render(content))
}

func SizedPanel(w, h int, content string) string {
	return lipgloss.NewStyle().
		Width(w).
		Height(h).
		Render(styles.Panel.Render(content))
}

func StatusBar(w int, content string) string {
	return lipgloss.NewStyle().
		Width(w).
		Background(theme.Bg).
		Foreground(theme.Text).
		Padding(0, 1).
		Render(content)
}

// Theme presets for different moods/contexts
var spaceTheme = StyleTheme{
	Primary:       lipgloss.Color("#6366F1"), // Indigo
	Secondary:     lipgloss.Color("#8B5CF6"), // Violet
	Accent:        lipgloss.Color("#06B6D4"), // Cyan
	Muted:         lipgloss.Color("#64748B"), // Slate
	Bg:            lipgloss.Color("#0F172A"), // Dark slate
	Text:          lipgloss.Color("#F1F5F9"), // Light slate
	Border:        lipgloss.Color("#334155"), // Slate
	Success:       lipgloss.Color("#10B981"), // Emerald
	Warning:       lipgloss.Color("#F59E0B"), // Amber
	Error:         lipgloss.Color("#EF4444"), // Red
	TextSecondary: lipgloss.Color("#94A3B8"), // Slate
}

var planetTheme = StyleTheme{
	Primary:       lipgloss.Color("#059669"), // Emerald
	Secondary:     lipgloss.Color("#0891B2"), // Cyan
	Accent:        lipgloss.Color("#7C3AED"), // Violet
	Muted:         lipgloss.Color("#6B7280"), // Gray
	Bg:            lipgloss.Color("#064E3B"), // Dark emerald
	Text:          lipgloss.Color("#ECFDF5"), // Light emerald
	Border:        lipgloss.Color("#047857"), // Emerald
	Success:       lipgloss.Color("#22C55E"), // Green
	Warning:       lipgloss.Color("#F59E0B"), // Amber
	Error:         lipgloss.Color("#DC2626"), // Red
	TextSecondary: lipgloss.Color("#9CA3AF"), // Gray
}

var dangerTheme = StyleTheme{
	Primary:       lipgloss.Color("#DC2626"), // Red
	Secondary:     lipgloss.Color("#EA580C"), // Orange
	Accent:        lipgloss.Color("#F59E0B"), // Amber
	Muted:         lipgloss.Color("#6B7280"), // Gray
	Bg:            lipgloss.Color("#7F1D1D"), // Dark red
	Text:          lipgloss.Color("#FEF2F2"), // Light red
	Border:        lipgloss.Color("#B91C1C"), // Red
	Success:       lipgloss.Color("#16A34A"), // Green
	Warning:       lipgloss.Color("#D97706"), // Amber
	Error:         lipgloss.Color("#B91C1C"), // Red
	TextSecondary: lipgloss.Color("#9CA3AF"), // Gray
}

// Theme switching functions
func UseSpaceTheme() {
	theme = spaceTheme
	rebuildStyles()
}

func UsePlanetTheme() {
	theme = planetTheme
	rebuildStyles()
}

func UseDangerTheme() {
	theme = dangerTheme
	rebuildStyles()
}

// Get available themes
func GetAvailableThemes() map[string]StyleTheme {
	return map[string]StyleTheme{
		"dark":   theme, // Current default
		"light":  lightTheme,
		"space":  spaceTheme,
		"planet": planetTheme,
		"danger": dangerTheme,
	}
}

// Apply theme by name
func ApplyTheme(name string) bool {
	themes := GetAvailableThemes()
	if selectedTheme, exists := themes[name]; exists {
		theme = selectedTheme
		rebuildStyles()
		return true
	}
	return false
}
