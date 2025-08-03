# Codebase Refactoring Analysis: Pain Points & Solutions

## Critical Issues Causing Movement Bug

### **1. Entity ID Hardcoding (URGENT)**
**Problem**: Systems are hardcoded to specific entity IDs:
```go
// SurfaceMovement.go line 102, 106
in, ok := ecs.Get[components.Input](w, 2)  // Hardcoded entity ID 2
p, ok := ecs.Get[components.Position](w, 2) // Hardcoded entity ID 2

// UI model creates player as entity ID varies
m.player = w.Create()  // Could be any ID, not necessarily 2
```

**Impact**: Movement fails because systems look for player at entity ID 2, but player entity might be ID 1, 3, etc.

**Solution**: Use component queries instead of hardcoded IDs:
```go
// Instead of hardcoded entity 2
ecs.View2Of[components.Input, components.Position](w).Each(func(t) {
    // Process all entities with Input + Position
})

// Or add Player component to identify the player entity
type Player struct{}
```

### **2. Duplicate & Conflicting Movement Systems**
**Problem**: Multiple movement systems with conflicting logic:
- `InputSystem` → sets `components.Velocity`
- `SpaceMovement` → reads `systems.Velocity` (different type!)
- `SurfaceMovement` → directly modifies Position, ignores Velocity
- `Movement` → reads `components.Velocity` but not used

**Impact**: Input creates one type of velocity, movement systems read different types.

**Solution**: Unify movement architecture around single velocity type.

### **3. Type Confusion**
**Problem**: Two different Velocity types:
```go
// pkg/components/velocity.go
type Velocity struct { VX, VY float64 }

// pkg/systems/space.go  
type Velocity struct { X, Y float64 }  // Different field names!
```

**Impact**: Systems can't communicate properly.

### **4. Component/System Boundary Violations**
**Problem**: Systems defining components:
```go
// In systems/space.go
type FuelTank struct{ Current int }  // Should be in components/
type Velocity struct{ X, Y float64 } // Conflicts with components.Velocity
```

**Impact**: Unclear ownership, import cycles, type conflicts.

## Major Architecture Pain Points

### **5. Missing Player Entity Management**
**Current**: Player entity ID varies, systems guess at ID 2
**Needed**: 
```go
type PlayerTag struct{}  // Mark the player entity
type WorldState struct {
    PlayerEntity ecs.Entity  // Track player globally
}
```

### **6. Layer System Over-Complexity**
**Problem**: Each system manually checks layer context:
```go
if ctx.CurrentLayer != ecs.LayerSpace { return }
```

**Impact**: Brittle, repetitive, easy to forget checks.

**Solution**: System registry handles layer activation automatically.

### **7. Inconsistent Input Handling**
**Problem**: Different input paradigms per layer:
- Space: Input → Velocity → Position (momentum-based physics)
- Surface: Input → Direct Position changes (grid-based movement)

**Design**: These systems are intentionally separate - spaceship movement and surface movement use different maps and mechanics. Space uses fuel/oxygen/warp cores, surface uses food/hunger. Layer transitions show full-screen loading ("Arriving on planet Toft") and pause/resume appropriate systems.

### **8. Monolithic UI Model**
**Problem**: UI model does too much:
- ECS world management
- Input translation  
- Rendering coordination
- Save/load
- System registry setup

**Solution**: Separate concerns into focused modules.

## File Organization Issues

### **9. Mixed Concerns in Files**
```
pkg/systems/space.go contains:
- FuelTank component (should be in components/)
- Velocity component (conflicts with components/)
- Multiple unrelated systems
- Helper functions
```

### **10. Missing Core Abstractions**
- No Player management system
- No unified movement interface
- No input event system
- No layer transition management

## Immediate Fixes Needed (Priority Order)

### **Priority 1: Fix Movement Bug**
1. **Add Player component** to identify player entity
2. **Fix SurfaceMovement** to use player component instead of hardcoded ID 2
3. **Unify Velocity types** - use components.Velocity everywhere
4. **Remove duplicate systems** - choose one movement approach

### **Priority 2: Clean Component/System Boundaries**
1. **Move components** from systems/ to components/
2. **Remove type conflicts** (duplicate Velocity types)
3. **Establish clear ownership** of each component type

### **Priority 3: Simplify Layer System**
1. **Remove manual layer checks** from individual systems
2. **Let SystemRegistry handle** layer-based activation
3. **Unified input processing** across all layers

### **Priority 4: Refactor Entity Management**
1. **Add Player tracking** to world state
2. **Component-based queries** instead of hardcoded entity IDs
3. **Centralized entity lifecycle** management

## Proposed Quick Fix for Movement

```go
// 1. Add to components/player.go
type Player struct{}

// 2. Fix SurfaceMovement to use component query
func (s SurfaceMovement) Update(dt float64, w *ecs.World) {
    ctx := ecs.GetWorldContext(w)
    if ctx.CurrentLayer != ecs.LayerPlanetSurface {
        return
    }
    
    // Find player entity via component query
    ecs.View3Of[Player, components.Input, components.Position](w).Each(func(t ecs.Tuple3[Player, components.Input, components.Position]) {
        dx, dy := 0.0, 0.0
        if t.B.Left { dx = -1 }
        if t.B.Right { dx = 1 }
        if t.B.Up { dy = -1 }
        if t.B.Down { dy = 1 }
        
        if dx != 0 || dy != 0 {
            t.C.X += dx
            t.C.Y += dy
            ecs.Add(w, t.E, *t.C)
        }
    })
}

// 3. Add Player component in UI model creation
ecs.Add(w, m.player, components.Player{})
```

## Long-term Architecture Goals

1. **Separate Movement Systems**: Maintain distinct space (momentum/physics) and surface (grid-based) movement with different resource systems
2. **Event-Driven Input**: Input events rather than direct component manipulation  
3. **Layer Abstraction**: Systems register for layers, automatic activation with pause/resume on transitions
4. **Component Purity**: All game state in components/, systems only contain logic
5. **Entity Services**: Centralized player, camera, world management
6. **Clean Separation**: UI ↔ Game Logic ↔ ECS boundaries
7. **Performance Target**: Game engine at 20 FPS, Bubble Tea UI at 60 FPS for smooth orbital mechanics and beautiful space phenomena

## Game Design Specifications

- **Layer Separation**: Space and surface are completely different maps with distinct mechanics
- **Resource Systems**: Space (fuel/oxygen/warp cores) vs Surface (food/hunger)
- **Transitions**: Full-screen loading with system pause ("Arriving on planet Toft")
- **Strict ECS**: Terminal is view layer only, all game logic in ECS
- **Single-player**: No multiplayer beyond seed sharing
- **Save System**: Quit-to-save with selectable save files
- **Performance**: 20 FPS game engine enables real-time orbital mechanics

This refactoring will eliminate the movement bug and create a maintainable codebase for the multi-planet roguelike vision.