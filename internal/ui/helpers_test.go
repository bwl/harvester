package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss/v2"
	"harvester/pkg/ecs"
)

func TestStyleBuilder(t *testing.T) {
	builder := NewStyleBuilder()
	if builder == nil {
		t.Fatal("NewStyleBuilder should return a valid builder")
	}

	// Test method chaining
	result := builder.
		Width(50).
		Height(10).
		Bold(true).
		Theme(ThemePrimary).
		Render("test content")

	if result == "" {
		t.Error("StyleBuilder should produce styled output")
	}
}

func TestStyleBuilderMethods(t *testing.T) {
	builder := NewStyleBuilder()

	// Test individual methods
	builder.Width(100)
	builder.Height(50)
	builder.PaddingHorizontal(2)
	builder.PaddingVertical(1)
	builder.Margin(1, 2, 3, 4)
	builder.Bold(true)
	builder.Italic(true)
	builder.Underline(true)
	builder.Border(lipgloss.RoundedBorder())
	builder.Align(lipgloss.Center)

	style := builder.Build()
	// Test that the style can render content
	result := style.Render("test")
	if result == "" {
		t.Error("StyleBuilder should build a valid style that can render content")
	}
}

func TestThemeColors(t *testing.T) {
	builder := NewStyleBuilder()

	colors := []ThemeColor{
		ThemePrimary,
		ThemeSecondary,
		ThemeAccent,
		ThemeMuted,
		ThemeSuccess,
		ThemeWarning,
		ThemeError,
	}

	for _, color := range colors {
		result := builder.Theme(color).Render("test")
		if result == "" {
			t.Errorf("Theme color %v should produce output", color)
		}
	}
}

func TestConditionalStyle(t *testing.T) {
	tests := []struct {
		state   GameState
		content string
	}{
		{StateNormal, "normal content"},
		{StateDanger, "danger content"},
		{StateWarning, "warning content"},
		{StateSuccess, "success content"},
		{StatePaused, "paused content"},
	}

	for _, test := range tests {
		result := ConditionalStyle(test.state, test.content)
		if result == "" {
			t.Errorf("ConditionalStyle(%v, %s) should return styled content", test.state, test.content)
		}
	}
}

func TestAnimationHelpers(t *testing.T) {
	// Test BlinkingText
	text := "blink me"
	result1 := BlinkingText(text, 10) // Should show text
	_ = BlinkingText(text, 40)        // Should show spaces (result not used in test)

	if result1 == "" {
		t.Error("BlinkingText should return content at some frames")
	}

	// Test FadingText
	fadeResult1 := FadingText(text, 1.0) // Full intensity
	fadeResult2 := FadingText(text, 0.5) // Half intensity
	fadeResult3 := FadingText(text, 0.0) // No intensity

	if fadeResult1 == "" || fadeResult2 == "" {
		t.Error("FadingText should return content at positive intensities")
	}

	if fadeResult3 != strings.Repeat(" ", len(text)) {
		t.Error("FadingText should return spaces at zero intensity")
	}
}

func TestComponentBuilder(t *testing.T) {
	builder := NewComponentBuilder()
	if builder == nil {
		t.Fatal("NewComponentBuilder should return a valid builder")
	}

	result := builder.
		Add("Component 1").
		Add("Component 2").
		Separator(" | ").
		Layout(lipgloss.Left).
		Build()

	if result == "" {
		t.Error("ComponentBuilder should produce output")
	}

	if !strings.Contains(result, "Component 1") || !strings.Contains(result, "Component 2") {
		t.Error("ComponentBuilder should include all added components")
	}
}

func TestComponentBuilderConditional(t *testing.T) {
	builder := NewComponentBuilder()

	result := builder.
		Add("Always shown").
		AddIf(true, "Conditionally shown").
		AddIf(false, "Never shown").
		Build()

	if !strings.Contains(result, "Always shown") {
		t.Error("ComponentBuilder should include unconditional components")
	}

	if !strings.Contains(result, "Conditionally shown") {
		t.Error("ComponentBuilder should include conditional components when true")
	}

	if strings.Contains(result, "Never shown") {
		t.Error("ComponentBuilder should not include conditional components when false")
	}
}

func TestComponentBuilderLayouts(t *testing.T) {
	builder := NewComponentBuilder()

	// Test vertical layout
	verticalResult := builder.
		Add("Top").
		Add("Bottom").
		Layout(lipgloss.Top).
		Build()

	if verticalResult == "" {
		t.Error("Vertical layout should produce output")
	}

	// Test horizontal layout
	builder = NewComponentBuilder()
	horizontalResult := builder.
		Add("Left").
		Add("Right").
		Layout(lipgloss.Left).
		Build()

	if horizontalResult == "" {
		t.Error("Horizontal layout should produce output")
	}
}

func TestComponentBuilderHelpers(t *testing.T) {
	builder := NewComponentBuilder()

	result := builder.
		Header("Test Header").
		Content("Test Content", Highlight).
		Footer("Test Footer").
		Build()

	if result == "" {
		t.Error("ComponentBuilder helpers should produce output")
	}
}

func TestGetStateColor(t *testing.T) {
	tests := []struct {
		layer    ecs.GameLayer
		expected ThemeColor
	}{
		{ecs.LayerSpace, ThemePrimary},
		{ecs.LayerPlanetSurface, ThemeSecondary},
		{ecs.LayerPlanetDeep, ThemeAccent},
	}

	for _, test := range tests {
		result := GetStateColor(test.layer)
		if result != test.expected {
			t.Errorf("GetStateColor(%v) = %v, expected %v", test.layer, result, test.expected)
		}
	}
}

func TestDynamicPanel(t *testing.T) {
	options := PanelOptions{
		Width:  50,
		Height: 20,
		Border: true,
		State:  StateNormal,
	}

	result := DynamicPanel("Test Title", "Test Content", StateNormal, options)
	if result == "" {
		t.Error("DynamicPanel should produce output")
	}
}

func TestStatWithTrend(t *testing.T) {
	tests := []struct {
		trend    TrendDirection
		expected string // We'll check if trend icon is included
	}{
		{TrendUp, "↗"},
		{TrendDown, "↘"},
		{TrendFlat, "→"},
	}

	for _, test := range tests {
		result := StatWithTrend("Test", "100", test.trend, StatGood)
		if result == "" {
			t.Errorf("StatWithTrend should produce output for trend %v", test.trend)
		}
		if !strings.Contains(result, test.expected) {
			t.Errorf("StatWithTrend should include trend icon %s", test.expected)
		}
	}
}

func TestTrendColorMapping(t *testing.T) {
	// Test internal function through StatWithTrend
	upResult := StatWithTrend("Test", "100", TrendUp, StatGood)
	downResult := StatWithTrend("Test", "100", TrendDown, StatGood)
	flatResult := StatWithTrend("Test", "100", TrendFlat, StatGood)

	// These should all produce different styled output
	if upResult == "" || downResult == "" || flatResult == "" {
		t.Error("All trend directions should produce output")
	}
}

func BenchmarkStyleBuilder(b *testing.B) {
	b.Run("BasicStyle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NewStyleBuilder().
				Width(50).
				Theme(ThemePrimary).
				Render("content")
		}
	})

	b.Run("ComplexStyle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NewStyleBuilder().
				Width(100).
				Height(50).
				PaddingHorizontal(2).
				Border(lipgloss.RoundedBorder()).
				Theme(ThemePrimary).
				Bold(true).
				Render("complex content")
		}
	})
}

func BenchmarkComponentBuilder(b *testing.B) {
	b.Run("SimpleComponent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NewComponentBuilder().
				Add("Component").
				Build()
		}
	})

	b.Run("ComplexComponent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NewComponentBuilder().
				Header("Header").
				Content("Content", Highlight).
				Footer("Footer").
				Separator("\n").
				Layout(lipgloss.Top).
				Build()
		}
	})
}

func BenchmarkAnimations(b *testing.B) {
	text := "Animation test content"

	b.Run("BlinkingText", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			BlinkingText(text, i%120) // Simulate frame counter
		}
	})

	b.Run("FadingText", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			intensity := float64(i%100) / 100.0
			FadingText(text, intensity)
		}
	})
}
