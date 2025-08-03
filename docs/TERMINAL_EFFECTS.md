# Terminal Game Effects Guide

This document explores cool visual effects and game-like interfaces achievable with Bubble Tea and the Charm ecosystem for Harvest of Stars.

## Core Charm Ecosystem Libraries

### **Bubble Tea** - TUI Framework
- MVU (Model-View-Update) architecture
- Event-driven input handling
- Full-screen or inline terminal apps
- Smooth window resizing and responsive layouts

### **Lip Gloss** - Styling & Layout
- Flexbox-style terminal layouts
- Rich color support (ANSI, hex, adaptive)
- Borders, padding, margins, alignment
- Gradients and advanced styling

### **Bubbles** - Component Library
- Pre-built interactive components (lists, tables, progress bars)
- Consistent styling and behavior
- Easy integration with Bubble Tea models

### **Harmonica** - Physics Animations
- Spring-based physics animations
- Three damping modes: under-damped (bouncy), critically-damped (smooth), over-damped (slow)
- Framework-agnostic 2D/3D animation support
- Natural motion for UI transitions

## Visual Effects for Harvest of Stars

### **Universe Expansion Animation**
```go
// Use Harmonica to animate map expansion
expandAnimation := harmonica.NewSpring(
    harmonica.FPS(60),
    harmonica.AngularVelocity(2.0),  // Speed of expansion
    harmonica.DampingRatio(0.8),     // Critically damped for smooth growth
)

// Animate viewport scaling as universe expands
viewportScale.Update(expandAnimation.Update(targetScale))
```

### **Galaxy Particle Effects**
```go
// Twinkling stars/galaxies using ANSI colors
galaxyStyles := []lipgloss.Style{
    lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Render, // Yellow
    lipgloss.NewStyle().Foreground(lipgloss.Color("45")).Render,  // Cyan
    lipgloss.NewStyle().Foreground(lipgloss.Color("201")).Render, // Pink
}

// Cycle through colors for twinkling effect
currentGalaxyStyle := galaxyStyles[tick%len(galaxyStyles)]
```

### **Progressive UI Reveal**
```go
// Status panels slide in with Harmonica spring animations
statusPanelAnimation := harmonica.NewSpring(
    harmonica.FPS(30),
    harmonica.AngularVelocity(1.5),
    harmonica.DampingRatio(0.6), // Under-damped for slight bounce
)

// Panel slides in from right edge
panelX.Update(statusPanelAnimation.Update(targetX))
```

### **Energy/Resource Bars**
```go
// Animated progress bars using Bubbles progress component
fuelBar := progress.New(progress.WithDefaultGradient())
hullBar := progress.New(progress.WithGradient("#FF0000", "#FF7F00"))

// Smooth value transitions with Harmonica
fuelAnimation.Update(currentFuel / maxFuel)
```

## Advanced Terminal Graphics

### **Box Drawing Characters**
```
┌─────────────────────────────────┬──────────────┐
│ EPOCH: 1  FUEL: ████████░░  85% │              │
├─────────────────────────────────┤   GALAXY     │
│                                 │   LOG        │
│           ∘ ∘ ∘ ∘ ∘             │              │
│         ∘ ∘ @ ∘ ∘ ∘             │              │
│           ∘ ∘ ★ ∘ ∘             │              │
│                                 │              │
├─────────────────────────────────┤              │
│ [hjkl] move [g] harvest [u] up  │              │
└─────────────────────────────────┴──────────────┘
```

### **Unicode Symbols for Game Elements**
- Player: `@` `⚫` `◉` `●`
- Galaxies: `★` `✦` `✧` `⭐` `*`
- Nebulae: `～` `≈` `∼` `⋈`
- Black Holes: `◯` `○` `⚬` `◦`
- Wormholes: `◈` `◇` `◆` `⬢`
- Energy: `⚡` `☀` `◉` `●`
- Void: `·` `∘` ` ` (space)

### **Color Schemes**
```go
// Deep space theme
var DeepSpaceTheme = struct {
    Void      lipgloss.Color
    Space     lipgloss.Color
    Galaxy    lipgloss.Color
    Player    lipgloss.Color
    Nebula    lipgloss.Color
    BlackHole lipgloss.Color
    UI        lipgloss.Color
}{
    Void:      lipgloss.Color("232"), // Almost black
    Space:     lipgloss.Color("240"), // Dark gray
    Galaxy:    lipgloss.Color("220"), // Bright yellow
    Player:    lipgloss.Color("45"),  // Cyan
    Nebula:    lipgloss.Color("135"), // Purple
    BlackHole: lipgloss.Color("88"),  // Dark red
    UI:        lipgloss.Color("250"), // Light gray
}

// Adaptive colors for light/dark terminals
adaptiveGalaxy := lipgloss.AdaptiveColor{
    Light: "#FFD700", // Gold in light mode
    Dark:  "#FFFF00", // Yellow in dark mode
}
```

## Game-Like Interface Elements

### **Smooth Scrolling Viewport**
```go
// Camera follows player with spring animation
cameraSpring := harmonica.NewSpring(
    harmonica.FPS(60),
    harmonica.AngularVelocity(3.0),   // Fast follow
    harmonica.DampingRatio(0.9),      // Minimal overshoot
)

// Update camera position toward player
cameraX.Update(cameraSpring.Update(playerX))
cameraY.Update(cameraSpring.Update(playerY))
```

### **Modal Upgrade Screens**
```go
// Overlay modal using Bubbles table component
upgradeTable := table.New(
    table.WithColumns([]table.Column{
        {Title: "Upgrade", Width: 15},
        {Title: "Level", Width: 6},
        {Title: "Cost", Width: 8},
    }),
    table.WithHeight(10),
    table.WithFocused(true),
)

// Animate modal slide-in
modalAnimation := harmonica.NewSpring(
    harmonica.FPS(30),
    harmonica.AngularVelocity(2.0),
    harmonica.DampingRatio(0.7),
)
```

### **Particle Systems**
```go
// Harvest animation - particles flying toward UI
type HarvestParticle struct {
    StartX, StartY float64
    TargetX, TargetY float64
    Progress float64
    Symbol rune
    Color lipgloss.Color
}

// Update particles with easing
for i := range particles {
    particles[i].Progress += deltaTime * 2.0
    if particles[i].Progress >= 1.0 {
        // Particle reached target, add to resources
        removeParticle(i)
    }
}
```

### **Dynamic Background Effects**
```go
// Cosmic background with moving stars
type StarField struct {
    Stars []Star
    Speed float64
}

type Star struct {
    X, Y float64
    Brightness float64
    Twinkle float64
}

// Update star twinkling
for i := range starField.Stars {
    star := &starField.Stars[i]
    star.Twinkle += deltaTime * 3.0
    star.Brightness = 0.5 + 0.5*math.Sin(star.Twinkle)
}
```

## Performance Optimizations

### **Viewport Culling**
```go
// Only render entities within camera bounds
visibleEntities := world.Query().
    Filter(func(pos components.Position) bool {
        return pos.X >= cameraX-margin && pos.X <= cameraX+viewWidth+margin &&
               pos.Y >= cameraY-margin && pos.Y <= cameraY+viewHeight+margin
    })
```

### **Dirty Region Updates**
```go
// Only re-render changed screen regions
type DirtyRegion struct {
    X, Y, Width, Height int
    NeedsUpdate bool
}

// Mark regions dirty when entities move
markDirty(oldPos.X, oldPos.Y, 1, 1)
markDirty(newPos.X, newPos.Y, 1, 1)
```

## Audio-Visual Feedback

### **Screen Shake Effects**
```go
// Camera shake on universe expansion
shakeIntensity := 2.0
shakeDuration := 0.3

if expanding {
    shakeX := (rand.Float64() - 0.5) * shakeIntensity
    shakeY := (rand.Float64() - 0.5) * shakeIntensity
    cameraX += int(shakeX)
    cameraY += int(shakeY)
}
```

### **Color Pulsing**
```go
// Pulse colors for important events
pulseBrightness := 0.5 + 0.5*math.Sin(time*4.0)
pulseColor := lipgloss.Color(fmt.Sprintf("#%02X%02X00", 
    int(255*pulseBrightness), int(255*pulseBrightness)))
```

## Recording & Documentation

### **VHS Integration**
Create beautiful demo recordings:
```tape
# demo.tape
Output demo.gif
Set Width 120
Set Height 40
Set Theme "Cosmic"

Type "go run ./cmd/game"
Enter
Sleep 2s
Type "hjkl" # Show movement
Sleep 1s
Type "g" # Harvest galaxy
Sleep 2s
```

Run with: `vhs demo.tape`

## Implementation Strategy

1. **Start with basic Lip Gloss styling** for immediate visual improvement
2. **Add Bubbles components** for professional UI elements
3. **Integrate Harmonica** for smooth animations
4. **Layer in particle effects** and advanced graphics
5. **Record with VHS** for documentation and showcasing

This creates a terminal experience that rivals graphical games while maintaining the charm and accessibility of text-based interfaces for Harvest of Stars.