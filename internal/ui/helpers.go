package ui

import (
	"strings"
	
	"github.com/charmbracelet/lipgloss"
	"harvester/pkg/ecs"
)

// StyleBuilder provides a fluent interface for building complex styles
type StyleBuilder struct {
	style lipgloss.Style
}

// NewStyleBuilder creates a new style builder
func NewStyleBuilder() *StyleBuilder {
	return &StyleBuilder{
		style: lipgloss.NewStyle(),
	}
}

// Width sets the width
func (sb *StyleBuilder) Width(w int) *StyleBuilder {
	sb.style = sb.style.Width(w)
	return sb
}

// Height sets the height
func (sb *StyleBuilder) Height(h int) *StyleBuilder {
	sb.style = sb.style.Height(h)
	return sb
}

// Padding sets padding
func (sb *StyleBuilder) Padding(top, right, bottom, left int) *StyleBuilder {
	sb.style = sb.style.Padding(top, right, bottom, left)
	return sb
}

// PaddingHorizontal sets horizontal padding
func (sb *StyleBuilder) PaddingHorizontal(horizontal int) *StyleBuilder {
	sb.style = sb.style.PaddingLeft(horizontal).PaddingRight(horizontal)
	return sb
}

// PaddingVertical sets vertical padding
func (sb *StyleBuilder) PaddingVertical(vertical int) *StyleBuilder {
	sb.style = sb.style.PaddingTop(vertical).PaddingBottom(vertical)
	return sb
}

// Margin sets margin
func (sb *StyleBuilder) Margin(top, right, bottom, left int) *StyleBuilder {
	sb.style = sb.style.Margin(top, right, bottom, left)
	return sb
}

// Background sets background color
func (sb *StyleBuilder) Background(color lipgloss.Color) *StyleBuilder {
	sb.style = sb.style.Background(color)
	return sb
}

// Foreground sets foreground color
func (sb *StyleBuilder) Foreground(color lipgloss.Color) *StyleBuilder {
	sb.style = sb.style.Foreground(color)
	return sb
}

// Bold sets bold styling
func (sb *StyleBuilder) Bold(bold bool) *StyleBuilder {
	sb.style = sb.style.Bold(bold)
	return sb
}

// Italic sets italic styling
func (sb *StyleBuilder) Italic(italic bool) *StyleBuilder {
	sb.style = sb.style.Italic(italic)
	return sb
}

// Underline sets underline styling
func (sb *StyleBuilder) Underline(underline bool) *StyleBuilder {
	sb.style = sb.style.Underline(underline)
	return sb
}

// Border sets border style
func (sb *StyleBuilder) Border(border lipgloss.Border) *StyleBuilder {
	sb.style = sb.style.Border(border)
	return sb
}

// BorderColor sets border color
func (sb *StyleBuilder) BorderColor(color lipgloss.Color) *StyleBuilder {
	sb.style = sb.style.BorderForeground(color)
	return sb
}

// Align sets text alignment
func (sb *StyleBuilder) Align(align lipgloss.Position) *StyleBuilder {
	sb.style = sb.style.Align(align)
	return sb
}

// Theme applies a theme-based color
func (sb *StyleBuilder) Theme(themeColor ThemeColor) *StyleBuilder {
	currentTheme := GetCurrentTheme()
	switch themeColor {
	case ThemePrimary:
		sb.style = sb.style.Foreground(currentTheme.Primary)
	case ThemeSecondary:
		sb.style = sb.style.Foreground(currentTheme.Secondary)
	case ThemeAccent:
		sb.style = sb.style.Foreground(currentTheme.Accent)
	case ThemeMuted:
		sb.style = sb.style.Foreground(currentTheme.Muted)
	case ThemeSuccess:
		sb.style = sb.style.Foreground(currentTheme.Success)
	case ThemeWarning:
		sb.style = sb.style.Foreground(currentTheme.Warning)
	case ThemeError:
		sb.style = sb.style.Foreground(currentTheme.Error)
	}
	return sb
}

// Build returns the final style
func (sb *StyleBuilder) Build() lipgloss.Style {
	return sb.style
}

// Render renders content with the built style
func (sb *StyleBuilder) Render(content string) string {
	return sb.style.Render(content)
}

// ThemeColor represents theme-based colors
type ThemeColor int

const (
	ThemePrimary ThemeColor = iota
	ThemeSecondary
	ThemeAccent
	ThemeMuted
	ThemeSuccess
	ThemeWarning
	ThemeError
)

// Conditional styling based on game state
type GameState int

const (
	StateNormal GameState = iota
	StateDanger
	StateWarning
	StateSuccess
	StatePaused
)

// ConditionalStyle applies styling based on game state
func ConditionalStyle(state GameState, content string) string {
	builder := NewStyleBuilder()
	
	switch state {
	case StateDanger:
		return builder.Theme(ThemeError).Bold(true).Render(content)
	case StateWarning:
		return builder.Theme(ThemeWarning).Render(content)
	case StateSuccess:
		return builder.Theme(ThemeSuccess).Render(content)
	case StatePaused:
		return builder.Theme(ThemeMuted).Italic(true).Render(content)
	default:
		return content
	}
}

// AnimatedText provides text animation helpers
type AnimationState struct {
	Frame      int
	MaxFrames  int
	Direction  int // 1 for forward, -1 for backward
}

// BlinkingText creates blinking effect
func BlinkingText(text string, frame int) string {
	if frame%60 < 30 { // Blink every 30 frames (at 60fps)
		return text
	}
	return strings.Repeat(" ", len(text))
}

// FadingText creates fading effect with different opacity levels
func FadingText(text string, intensity float64) string {
	if intensity <= 0 {
		return strings.Repeat(" ", len(text))
	}
	if intensity >= 1 {
		return text
	}
	
	// Simulate fading with different colors
	builder := NewStyleBuilder()
	if intensity > 0.7 {
		return builder.Theme(ThemePrimary).Render(text)
	} else if intensity > 0.4 {
		return builder.Theme(ThemeMuted).Render(text)
	} else {
		return builder.Theme(ThemeMuted).Foreground(lipgloss.Color("237")).Render(text)
	}
}

// Component composition helpers
type ComponentBuilder struct {
	components []string
	separator  string
	layout     lipgloss.Position
}

// NewComponentBuilder creates a new component builder
func NewComponentBuilder() *ComponentBuilder {
	return &ComponentBuilder{
		components: make([]string, 0),
		separator:  "",
		layout:     lipgloss.Left,
	}
}

// Add adds a component
func (cb *ComponentBuilder) Add(component string) *ComponentBuilder {
	cb.components = append(cb.components, component)
	return cb
}

// AddIf adds a component conditionally
func (cb *ComponentBuilder) AddIf(condition bool, component string) *ComponentBuilder {
	if condition {
		cb.components = append(cb.components, component)
	}
	return cb
}

// Separator sets the separator between components
func (cb *ComponentBuilder) Separator(sep string) *ComponentBuilder {
	cb.separator = sep
	return cb
}

// Layout sets the layout direction
func (cb *ComponentBuilder) Layout(layout lipgloss.Position) *ComponentBuilder {
	cb.layout = layout
	return cb
}

// Header adds a header component
func (cb *ComponentBuilder) Header(text string) *ComponentBuilder {
	return cb.Add(Header(text))
}

// Content adds content with optional styling
func (cb *ComponentBuilder) Content(text string, style func(string) string) *ComponentBuilder {
	if style != nil {
		return cb.Add(style(text))
	}
	return cb.Add(text)
}

// Footer adds a footer component
func (cb *ComponentBuilder) Footer(text string) *ComponentBuilder {
	return cb.Add(Muted(text))
}

// Build returns the final composed component
func (cb *ComponentBuilder) Build() string {
	if len(cb.components) == 0 {
		return ""
	}
	
	if cb.layout == lipgloss.Top || cb.layout == lipgloss.Bottom {
		// Vertical layout
		if cb.separator != "" {
			var separated []string
			for i, comp := range cb.components {
				separated = append(separated, comp)
				if i < len(cb.components)-1 {
					separated = append(separated, cb.separator)
				}
			}
			return lipgloss.JoinVertical(cb.layout, separated...)
		}
		return lipgloss.JoinVertical(cb.layout, cb.components...)
	} else {
		// Horizontal layout
		if cb.separator != "" {
			var separated []string
			for i, comp := range cb.components {
				separated = append(separated, comp)
				if i < len(cb.components)-1 {
					separated = append(separated, cb.separator)
				}
			}
			return lipgloss.JoinHorizontal(cb.layout, separated...)
		}
		return lipgloss.JoinHorizontal(cb.layout, cb.components...)
	}
}

// State-aware color schemes
func GetStateColor(state ecs.GameLayer) ThemeColor {
	switch state {
	case ecs.LayerSpace:
		return ThemePrimary
	case ecs.LayerPlanetSurface:
		return ThemeSecondary
	case ecs.LayerPlanetDeep:
		return ThemeAccent
	default:
		return ThemeMuted
	}
}

// Dynamic panel styling based on content
func DynamicPanel(title, content string, state GameState, options PanelOptions) string {
	builder := NewComponentBuilder()
	
	// Add title with state-aware styling
	if title != "" {
		styledTitle := ConditionalStyle(state, title)
		builder.Header(styledTitle)
	}
	
	// Add content
	if content != "" {
		builder.Content(content, nil)
	}
	
	panel := builder.Layout(lipgloss.Top).Build()
	
	// Apply panel styling
	if options.Border {
		return Bordered(panel)
	}
	if options.Width > 0 || options.Height > 0 {
		return Sized(options.Width, options.Height, panel)
	}
	
	return panel
}

// PanelOptions for dynamic panel creation
type PanelOptions struct {
	Width  int
	Height int
	Border bool
	State  GameState
}

// Advanced stat display with trend indicators
func StatWithTrend(label, current string, trend TrendDirection, status StatStatus) string {
	var trendIcon string
	switch trend {
	case TrendUp:
		trendIcon = "↗"
	case TrendDown:
		trendIcon = "↘"
	case TrendFlat:
		trendIcon = "→"
	default:
		trendIcon = ""
	}
	
	stat := Stat(label, current, status)
	if trendIcon != "" {
		trendStyle := NewStyleBuilder().Theme(getTrendColor(trend)).Build()
		stat += " " + trendStyle.Render(trendIcon)
	}
	
	return stat
}

type TrendDirection int

const (
	TrendFlat TrendDirection = iota
	TrendUp
	TrendDown
)

func getTrendColor(trend TrendDirection) ThemeColor {
	switch trend {
	case TrendUp:
		return ThemeSuccess
	case TrendDown:
		return ThemeError
	default:
		return ThemeMuted
	}
}

