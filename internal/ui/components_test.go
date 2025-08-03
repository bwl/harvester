package ui

import (
	"strings"
	"testing"

	"harvester/pkg/ecs"
)

func TestStatusBarComponent(t *testing.T) {
	sections := []StatusSection{
		{Label: "Health: ", Value: "100", Style: nil},
		{Label: "Energy: ", Value: "75", Style: Highlight},
	}

	result := StatusBarComponent(80, sections)
	if result == "" {
		t.Error("StatusBarComponent should produce output")
	}

	if !strings.Contains(result, "Health") || !strings.Contains(result, "Energy") {
		t.Error("StatusBarComponent should include all sections")
	}
}

func TestQuestPanel(t *testing.T) {
	data := QuestPanelData{
		Status: "In Progress",
	}

	result := QuestPanel(data)
	if result == "" {
		t.Error("QuestPanel should produce output")
	}

	if !strings.Contains(result, "QUEST") {
		t.Error("QuestPanel should include quest header")
	}

	if !strings.Contains(result, "In Progress") {
		t.Error("QuestPanel should include status")
	}
}

func TestControlsPanel(t *testing.T) {
	groups := []ControlsGroup{
		{
			Title: "Movement",
			Items: []ControlItem{
				{"WASD", "move around"},
				{"Space", "jump"},
			},
		},
		{
			Title: "Actions",
			Items: []ControlItem{
				{"E", "interact"},
				{"Q", "quit"},
			},
		},
	}

	result := ControlsPanel(groups)
	if result == "" {
		t.Error("ControlsPanel should produce output")
	}

	if !strings.Contains(result, "Movement") || !strings.Contains(result, "Actions") {
		t.Error("ControlsPanel should include all groups")
	}

	if !strings.Contains(result, "WASD") || !strings.Contains(result, "interact") {
		t.Error("ControlsPanel should include control items")
	}
}

func TestMapPanel(t *testing.T) {
	content := "Map content here"

	// Test without border
	opts := MapPanelOptions{
		Width:  50,
		Height: 20,
		Border: false,
	}

	result := MapPanel(content, opts)
	if result == "" {
		t.Error("MapPanel should produce output")
	}

	// Test with border
	opts.Border = true
	result = MapPanel(content, opts)
	if result == "" {
		t.Error("MapPanel with border should produce output")
	}

	// Test with title and border
	opts.Title = "World Map"
	result = MapPanel(content, opts)
	if result == "" {
		t.Error("MapPanel with title should produce output")
	}

	if !strings.Contains(result, "World Map") {
		t.Error("MapPanel should include title when specified")
	}
}

func TestLogPanel(t *testing.T) {
	// Test empty log
	result := LogPanel([]LogMessage{}, 50)
	if result == "" {
		t.Error("LogPanel should handle empty messages")
	}

	// Test with messages
	messages := []LogMessage{
		{Text: "Info message", Type: LogInfo},
		{Text: "Warning message", Type: LogWarning},
		{Text: "Error message", Type: LogError},
		{Text: "Success message", Type: LogSuccess},
	}

	result = LogPanel(messages, 50)
	if result == "" {
		t.Error("LogPanel should produce output with messages")
	}

	if !strings.Contains(result, "Info message") {
		t.Error("LogPanel should include all message types")
	}
}

func TestPlayerStatsComponent(t *testing.T) {
	stats := PlayerStatsData{
		Fuel:  85,
		Hull:  60,
		Drive: 3,
	}

	result := PlayerStatsComponent(stats)
	if result == "" {
		t.Error("PlayerStatsComponent should produce output")
	}

	if !strings.Contains(result, "85") || !strings.Contains(result, "60") || !strings.Contains(result, "3") {
		t.Error("PlayerStatsComponent should include all stat values")
	}
}

func TestLocationComponent(t *testing.T) {
	location := LocationData{
		Layer:  "Space",
		Planet: 42,
		Depth:  5,
	}

	result := LocationComponent(location)
	if result == "" {
		t.Error("LocationComponent should produce output")
	}

	if !strings.Contains(result, "Space") || !strings.Contains(result, "42") || !strings.Contains(result, "5") {
		t.Error("LocationComponent should include all location data")
	}
}

func TestGameInfoComponent(t *testing.T) {
	info := GameInfoData{
		Tick: 12345,
	}

	result := GameInfoComponent(info)
	if result == "" {
		t.Error("GameInfoComponent should produce output")
	}

	if !strings.Contains(result, "12345") {
		t.Error("GameInfoComponent should include tick count")
	}
}

func TestEnhancedStatusBar(t *testing.T) {
	location := LocationData{Layer: "Space", Planet: 1, Depth: 0}
	stats := PlayerStatsData{Fuel: 100, Hull: 80, Drive: 2}
	info := GameInfoData{Tick: 500}

	result := EnhancedStatusBar(100, location, stats, info)
	if result == "" {
		t.Error("EnhancedStatusBar should produce output")
	}
}

func TestDynamicQuestPanel(t *testing.T) {
	data := QuestPanelData{Status: "Active"}

	result := DynamicQuestPanel(data, StateNormal)
	if result == "" {
		t.Error("DynamicQuestPanel should produce output")
	}

	// Test different states
	result = DynamicQuestPanel(data, StateDanger)
	if result == "" {
		t.Error("DynamicQuestPanel should handle danger state")
	}
}

func TestAnimatedStatusComponent(t *testing.T) {
	data := AnimatedStatusData{
		Label:    "Health ",
		Value:    "25",
		Status:   StatDanger,
		Trend:    TrendDown,
		Animated: true,
		Frame:    30,
	}

	result := AnimatedStatusComponent(data)
	if result == "" {
		t.Error("AnimatedStatusComponent should produce output")
	}

	// Test non-animated version
	data.Animated = false
	result = AnimatedStatusComponent(data)
	if result == "" {
		t.Error("AnimatedStatusComponent should handle non-animated state")
	}
}

func TestAdvancedPlayerStatsComponent(t *testing.T) {
	current := PlayerStatsData{Fuel: 80, Hull: 90, Drive: 2}
	previous := PlayerStatsData{Fuel: 85, Hull: 85, Drive: 2}

	result := AdvancedPlayerStatsComponent(current, previous, 60)
	if result == "" {
		t.Error("AdvancedPlayerStatsComponent should produce output")
	}

	// Should show trends based on comparison
	if !strings.Contains(result, "80") || !strings.Contains(result, "90") {
		t.Error("AdvancedPlayerStatsComponent should include current values")
	}
}

func TestThemedPanel(t *testing.T) {
	result := ThemedPanel("Test Panel", "Content", ecs.LayerSpace, 50, 20)
	if result == "" {
		t.Error("ThemedPanel should produce output")
	}

	// Test different layers
	result = ThemedPanel("Surface Panel", "Content", ecs.LayerPlanetSurface, 50, 20)
	if result == "" {
		t.Error("ThemedPanel should handle different layers")
	}
}

func TestResponsiveControlsPanel(t *testing.T) {
	groups := []ControlsGroup{
		{
			Title: "Movement",
			Items: []ControlItem{
				{"h j k l", "move"},
				{"shift+move", "run"},
			},
		},
		{
			Title: "Actions",
			Items: []ControlItem{
				{"> ", "enter"},
				{"q ", "quit"},
				{"i", "inventory"},
			},
		},
	}

	// Test full mode
	result := ResponsiveControlsPanel(groups, 30)
	if result == "" {
		t.Error("ResponsiveControlsPanel should produce output in full mode")
	}

	// Test compact mode
	compactResult := ResponsiveControlsPanel(groups, 20)
	if compactResult == "" {
		t.Error("ResponsiveControlsPanel should produce output in compact mode")
	}

	// Compact mode should show fewer controls (check by counting control items, not colons)
	fullControlCount := strings.Count(result, "move") + strings.Count(result, "enter") + strings.Count(result, "quit") + strings.Count(result, "inventory")
	compactControlCount := strings.Count(compactResult, "move") + strings.Count(compactResult, "enter") + strings.Count(compactResult, "quit") + strings.Count(compactResult, "inventory")

	// In compact mode, non-essential controls like "inventory" should be filtered out
	if compactControlCount >= fullControlCount {
		t.Error("Compact mode should show fewer controls than full mode")
	}
}

func TestEssentialControlFiltering(t *testing.T) {
	tests := []struct {
		key         string
		isEssential bool
	}{
		{"h j k l", true},
		{"> ", true},
		{"q ", true},
		{"i", false},
		{"shift+move", false},
		{"ctrl+s", false},
	}

	for _, test := range tests {
		result := isEssentialControl(test.key)
		if result != test.isEssential {
			t.Errorf("isEssentialControl(%s) = %v, expected %v", test.key, result, test.isEssential)
		}
	}
}

func BenchmarkComponents(b *testing.B) {
	location := LocationData{Layer: "Space", Planet: 1, Depth: 0}
	stats := PlayerStatsData{Fuel: 100, Hull: 80, Drive: 2}
	info := GameInfoData{Tick: 500}

	b.Run("LocationComponent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			LocationComponent(location)
		}
	})

	b.Run("PlayerStatsComponent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			PlayerStatsComponent(stats)
		}
	})

	b.Run("EnhancedStatusBar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			EnhancedStatusBar(100, location, stats, info)
		}
	})
}

func BenchmarkAdvancedComponents(b *testing.B) {
	current := PlayerStatsData{Fuel: 80, Hull: 90, Drive: 2}
	previous := PlayerStatsData{Fuel: 85, Hull: 85, Drive: 2}

	b.Run("AdvancedPlayerStatsComponent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			AdvancedPlayerStatsComponent(current, previous, i%120)
		}
	})

	b.Run("ThemedPanel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ThemedPanel("Test", "Content", ecs.LayerSpace, 50, 20)
		}
	})
}
