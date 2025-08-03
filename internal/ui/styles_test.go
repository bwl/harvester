package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestThemeCreation(t *testing.T) {
	// Test default theme
	currentTheme := GetCurrentTheme()
	if currentTheme.Primary == "" {
		t.Error("Default theme should have a primary color")
	}
	if currentTheme.Text == "" {
		t.Error("Default theme should have a text color")
	}
}

func TestThemeSwitching(t *testing.T) {
	// Store original theme
	originalTheme := GetCurrentTheme()
	
	// Test light theme
	UseLightTheme()
	lightTheme := GetCurrentTheme()
	if lightTheme.Bg != lipgloss.Color("255") {
		t.Error("Light theme should have white background")
	}
	
	// Test space theme
	UseSpaceTheme()
	spaceTheme := GetCurrentTheme()
	if spaceTheme.Primary != lipgloss.Color("#6366F1") {
		t.Error("Space theme should have indigo primary color")
	}
	
	// Test planet theme
	UsePlanetTheme()
	planetTheme := GetCurrentTheme()
	if planetTheme.Primary != lipgloss.Color("#059669") {
		t.Error("Planet theme should have emerald primary color")
	}
	
	// Test danger theme
	UseDangerTheme()
	dangerTheme := GetCurrentTheme()
	if dangerTheme.Primary != lipgloss.Color("#DC2626") {
		t.Error("Danger theme should have red primary color")
	}
	
	// Restore original theme
	SetCustomTheme(originalTheme)
}

func TestApplyThemeByName(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"light", true},
		{"dark", true},
		{"space", true},
		{"planet", true},
		{"danger", true},
		{"nonexistent", false},
		{"", false},
	}
	
	for _, test := range tests {
		result := ApplyTheme(test.name)
		if result != test.expected {
			t.Errorf("ApplyTheme(%s) = %v, expected %v", test.name, result, test.expected)
		}
	}
}

func TestStyleHelperFunctions(t *testing.T) {
	// Test basic helpers
	panelResult := Panel("test content")
	if panelResult == "" {
		t.Error("Panel helper should return styled content")
	}
	
	borderedResult := Bordered("test content")
	if borderedResult == "" {
		t.Error("Bordered helper should return styled content")
	}
	
	headerResult := Header("Test Header")
	if headerResult == "" {
		t.Error("Header helper should return styled content")
	}
	
	mutedResult := Muted("muted text")
	if mutedResult == "" {
		t.Error("Muted helper should return styled content")
	}
	
	highlightResult := Highlight("highlighted text")
	if highlightResult == "" {
		t.Error("Highlight helper should return styled content")
	}
}

func TestSizedHelper(t *testing.T) {
	result := Sized(10, 5, "content")
	if result == "" {
		t.Error("Sized helper should return styled content")
	}
	
	// Test with zero dimensions
	result = Sized(0, 0, "content")
	if result == "" {
		t.Error("Sized helper should handle zero dimensions")
	}
}

func TestStatFunction(t *testing.T) {
	tests := []struct {
		label  string
		value  string
		status StatStatus
	}{
		{"Health", "100", StatGood},
		{"Energy", "50", StatWarning},
		{"Hull", "10", StatDanger},
	}
	
	for _, test := range tests {
		result := Stat(test.label, test.value, test.status)
		if result == "" {
			t.Errorf("Stat(%s, %s, %v) should return styled content", test.label, test.value, test.status)
		}
	}
}

func TestGetStatColor(t *testing.T) {
	tests := []struct {
		current  int
		max      int
		expected StatStatus
	}{
		{100, 100, StatGood},   // 100%
		{80, 100, StatGood},    // 80%
		{50, 100, StatWarning}, // 50%
		{20, 100, StatDanger},  // 20%
		{0, 100, StatDanger},   // 0%
	}
	
	for _, test := range tests {
		result := GetStatColor(test.current, test.max)
		if result != test.expected {
			t.Errorf("GetStatColor(%d, %d) = %v, expected %v", test.current, test.max, result, test.expected)
		}
	}
}

func TestAdvancedHelpers(t *testing.T) {
	// Test PanelWithBorder
	result := PanelWithBorder("content")
	if result == "" {
		t.Error("PanelWithBorder should return styled content")
	}
	
	// Test SizedPanel
	result = SizedPanel(20, 10, "content")
	if result == "" {
		t.Error("SizedPanel should return styled content")
	}
	
	// Test StatusBar
	result = StatusBar(50, "status content")
	if result == "" {
		t.Error("StatusBar should return styled content")
	}
}

func TestThemeConsistency(t *testing.T) {
	themes := GetAvailableThemes()
	
	for name, theme := range themes {
		// Check that all required colors are defined
		if theme.Primary == "" {
			t.Errorf("Theme %s missing Primary color", name)
		}
		if theme.Secondary == "" {
			t.Errorf("Theme %s missing Secondary color", name)
		}
		if theme.Text == "" {
			t.Errorf("Theme %s missing Text color", name)
		}
		if theme.Bg == "" {
			t.Errorf("Theme %s missing Background color", name)
		}
		if theme.Success == "" {
			t.Errorf("Theme %s missing Success color", name)
		}
		if theme.Warning == "" {
			t.Errorf("Theme %s missing Warning color", name)
		}
		if theme.Error == "" {
			t.Errorf("Theme %s missing Error color", name)
		}
	}
}

func BenchmarkStyleHelpers(b *testing.B) {
	b.Run("Panel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Panel("test content")
		}
	})
	
	b.Run("Header", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Header("Test Header")
		}
	})
	
	b.Run("Stat", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Stat("Label", "Value", StatGood)
		}
	})
}

func BenchmarkThemeSwitching(b *testing.B) {
	b.Run("UseSpaceTheme", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			UseSpaceTheme()
		}
	})
	
	b.Run("ApplyTheme", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ApplyTheme("space")
		}
	})
}