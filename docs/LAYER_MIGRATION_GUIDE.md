# Layer-Based Rendering Migration Guide

This guide explains how to migrate from the old glyph-based rendering system to the new lipgloss v2 Layer-based system.

## Overview

The rendering system has been upgraded to use lipgloss v2's Canvas and Layer functionality, providing:

- **Simpler code**: No more manual glyph matrix management
- **Better performance**: Canvas handles composition optimally
- **Rich styling**: Full lipgloss styling capabilities
- **Automatic layout**: Built-in text wrapping, alignment, borders
- **Backward compatibility**: Old glyph-based content still works

## Architecture

### Old System (Glyph-based)
```go
type RenderableContent interface {
    GetLayer() Layer
    GetZ() int
    GetPosition() Position
    GetBounds() Bounds
    GetGlyphs() [][]Glyph  // Manual glyph matrix
}
```

### New System (Layer-based)
```go
type LayerContent interface {
    GetLayer() Layer
    GetZ() int
    ToLipglossLayer() *lipgloss.Layer  // Direct lipgloss integration
}
```

## Migration Examples

### 1. Simple Text Content

**Old way:**
```go
func (t *textContent) GetGlyphs() [][]rendering.Glyph {
    // Manually create glyph matrix
    glyphs := make([][]rendering.Glyph, height)
    // ... complex glyph generation code
    return glyphs
}
```

**New way:**
```go
content := rendering.NewStyledContent(
    rendering.LayerUI,
    rendering.ZUI,
    "Hello, World!",
).WithStyle(
    lipgloss.NewStyle().
        Foreground(lipgloss.Color("#ffffff")).
        Background(lipgloss.Color("#1a1a2e")).
        Padding(1, 2).
        Border(lipgloss.RoundedBorder()),
).WithPosition(10, 5)

renderer.RegisterLayerContent(content)
```

### 2. UI Panels

**Old way:**
```go
type uiPanel struct {
    content string
    w, h    int
    // ... manual styling and positioning
}

func (u *uiPanel) GetGlyphs() [][]rendering.Glyph {
    // Complex glyph matrix creation with manual styling
}
```

**New way:**
```go
panel := ui.NewLayerTextPanel(
    "Panel Title",
    "Panel content with automatic\nwrapping and styling",
    40,
).WithPosition(10, 5)

renderer.RegisterLayerContent(panel)
```

### 3. Status Bars

**Old way:**
```go
// Manual character-by-character rendering
func createStatusBar() [][]rendering.Glyph {
    // ... complex glyph assembly
}
```

**New way:**
```go
status := ui.NewLayerStatusBar(80).
    AddItem("Health", "100").
    AddItem("Score", "1337").
    AddItem("Level", "5").
    WithPosition(0, 19)

renderer.RegisterLayerContent(status)
```

### 4. Game Screens

**Old way:**
```go
func (s *StartScreen) RegisterContent(renderer *rendering.ViewRenderer) {
    // Register multiple glyph-based components
    renderer.RegisterContent(background)
    renderer.RegisterContent(title)
    renderer.RegisterContent(menu)
}
```

**New way:**
```go
func (s *LayerStartScreen) RegisterLayerContent(renderer *rendering.CanvasRenderer) {
    // Create styled components with fluent API
    title := rendering.NewStyledContent(
        rendering.LayerUI,
        rendering.ZUI,
        "🌌 HARVESTER",
    ).WithStyle(titleStyle).WithPosition(10, 3)
    
    renderer.RegisterLayerContent(title)
}
```

## Available Components

### Core Components
- `StyledContent`: Basic styled text content
- `CompositeContent`: Multiple layers composed together
- `RenderableContentAdapter`: Backward compatibility wrapper

### UI Components
- `LayerTVFrame`: Styled TV frame border
- `LayerBackground`: Gradient/pattern backgrounds
- `LayerTextPanel`: Text panels with borders and padding
- `LayerStatusBar`: Status bars with items
- `LayerGameContent`: Game entities with positioning

## Migration Steps

### 1. Update Renderer Registration
```go
// Old
func (s *Screen) RegisterContent(renderer *rendering.ViewRenderer) {
    renderer.RegisterContent(oldContent)
}

// New  
func (s *Screen) RegisterLayerContent(renderer *rendering.CanvasRenderer) {
    renderer.RegisterLayerContent(newContent)
}
```

### 2. Replace Glyph Generation
```go
// Old - Manual glyph matrix
func (c *Component) GetGlyphs() [][]rendering.Glyph {
    // 50+ lines of glyph generation
}

// New - Lipgloss styling
func (c *Component) ToLipglossLayer() *lipgloss.Layer {
    styledContent := lipgloss.NewStyle().
        Foreground(color).
        Background(bgColor).
        Render(content)
    
    return lipgloss.NewLayer(styledContent).X(x).Y(y).Z(z)
}
```

### 3. Use Fluent API
```go
// Chain styling operations
component := rendering.NewStyledContent(layer, z, text).
    WithForeground(lipgloss.Color("#fff")).
    WithBackground(lipgloss.Color("#333")).
    WithBorder(lipgloss.RoundedBorder()).
    WithPadding(1, 2, 1, 2).
    WithPosition(x, y)
```

## Best Practices

### 1. Layer Organization
```go
const (
    ZBackground = 0      // Backgrounds
    ZGame       = 50     // Game entities  
    ZContent    = 100    // Main content
    ZUI         = 500    // UI elements
    ZHUD        = 700    // HUD overlays
    ZMenu       = 900    // Menus
    ZFrame      = 1000   // Frames/borders
)
```

### 2. Component Composition
```go
// Create composite components
composite := rendering.NewCompositeContent(rendering.LayerUI, rendering.ZUI).
    AddChild(background).
    AddChild(border).
    AddChild(content)
```

### 3. Reusable Styles
```go
// Define common styles
var (
    PanelStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#3a3a4a")).
        Background(lipgloss.Color("#1a1a2e")).
        Padding(1, 2)
    
    TitleStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#ffd700")).
        Bold(true).
        Align(lipgloss.Center)
)
```

## Performance Notes

- **Canvas optimization**: lipgloss Canvas handles dirty regions automatically
- **Layer caching**: Layers are cached by lipgloss for better performance  
- **Reduced complexity**: No manual alpha blending or glyph management
- **Memory efficiency**: Less memory allocation than glyph matrices

## Backward Compatibility

The system maintains full backward compatibility:

```go
// Old RenderableContent still works
renderer.RegisterContent(oldGlyphContent)

// New LayerContent works alongside
renderer.RegisterLayerContent(newLayerContent)

// Both render together in the same canvas
output := renderer.Render()
```

## Common Patterns

### Dynamic Content Updates
```go
// Update content and re-register
status.items[0].Value = "95"  // Update health
renderer.RegisterLayerContent(status)  // Re-register
```

### Conditional Rendering
```go
if showHelp {
    help := ui.NewLayerTextPanel("Help", helpText, 40)
    renderer.RegisterLayerContent(help)
}
```

### Animation Support
```go
// Create animated content by updating position/style each frame
for frame := 0; frame < 60; frame++ {
    alpha := float64(frame) / 60.0
    style := baseStyle.Copy().Foreground(fadeColor(alpha))
    
    content := rendering.NewStyledContent(layer, z, text).
        WithStyle(style).
        WithPosition(x, y+frame)
    
    renderer.RegisterLayerContent(content)
}
```

This new system provides a much more maintainable and feature-rich approach to terminal UI rendering while maintaining full compatibility with existing code.