# Style Refactor Plan

This document outlines the plan to refactor the lipgloss styling system for better maintainability, reusability, and cleaner code organization.

## Current State

The current styling system in `internal/ui/model.go` has:
- All theme and style definitions mixed with UI logic
- Repetitive `style.Copy().Width().Render()` patterns
- Complex inline style building
- Limited reusability across components

## Refactoring Goals

1. **Better Organization** - Separate styling from UI logic
2. **Reduced Repetition** - Create reusable style helpers
3. **Improved Readability** - Cleaner, more semantic styling API
4. **Enhanced Maintainability** - Easier to modify themes and styles
5. **Component-Based Architecture** - Modular, testable render functions

## Phase 1: Style Extraction & Organization

### Task 1.1: Create Style Module
- [ ] Create `internal/ui/styles.go`
- [ ] Move `StyleTheme` struct and `theme` variable
- [ ] Move `styles` struct with all style definitions
- [ ] Export necessary styles for use in other UI files

### Task 1.2: Create Style Helpers
- [ ] `Panel(content string) string` - Standard panel with padding/background
- [ ] `Bordered(content string) string` - Content with rounded borders
- [ ] `Sized(w, h int, content string) string` - Content with fixed dimensions
- [ ] `Header(text string) string` - Consistent header styling
- [ ] `Stat(label, value string, status StatStatus) string` - Color-coded stat display

### Task 1.3: Theme Management
- [ ] Create theme switching capability (dark/light)
- [ ] Add theme validation
- [ ] Create theme builder functions

## Phase 2: Component Refactoring

### Task 2.1: Status Bar Component
- [ ] Extract `renderStatusBar` improvements
- [ ] Create `StatusSection` component for location/stats/tick
- [ ] Add responsive status bar layout
- [ ] Create stat color helper with enums

### Task 2.2: Panel Components
- [ ] Create `QuestPanel` component
- [ ] Create `ControlsPanel` component  
- [ ] Create `MapPanel` component with border options
- [ ] Create `LogPanel` component with message type styling

### Task 2.3: Layout System
- [ ] Create `Layout` struct for managing panel dimensions
- [ ] Add responsive layout calculations
- [ ] Create layout presets (full, compact, mobile)
- [ ] Add layout validation

## Phase 3: Advanced Style Patterns

### Task 3.1: Style Builder Pattern
```go
// Example API:
style := ui.NewStyleBuilder().
    Width(30).
    Height(10).
    Border(ui.RoundedBorder).
    Theme(ui.PrimaryTheme).
    Build()
```

### Task 3.2: Component Composition
```go
// Example API:
panel := ui.Panel().
    Header("QUEST").
    Content(questStatus).
    Footer(controls).
    Render()
```

### Task 3.3: Dynamic Styling
- [ ] Conditional styling based on game state
- [ ] Animation/transition helpers
- [ ] State-aware color schemes

## Phase 4: Testing & Polish

### Task 4.1: Style Testing
- [ ] Unit tests for style helpers
- [ ] Visual regression tests for components
- [ ] Theme switching tests

### Task 4.2: Performance Optimization
- [ ] Style caching for expensive operations
- [ ] Lazy style computation
- [ ] Memory usage optimization

### Task 4.3: Documentation
- [ ] Style guide documentation
- [ ] Component usage examples
- [ ] Theme customization guide

## File Structure (Post-Refactor)

```
internal/ui/
├── styles.go          # Theme and style definitions
├── components.go      # Reusable UI components
├── layout.go          # Layout calculation helpers
├── model.go           # Main UI model (simplified)
└── helpers.go         # Style helper functions
```

## Breaking Changes

- Move from inline styling to component-based approach
- Centralize theme configuration
- Change some function signatures for consistency

## Benefits

1. **Maintainability** - Easier to modify and extend styling
2. **Consistency** - Unified approach to styling across components
3. **Reusability** - Components can be used in multiple contexts
4. **Testability** - Individual components can be unit tested
5. **Performance** - Optimized style computation and caching
6. **Developer Experience** - Cleaner, more intuitive API

## Migration Strategy

1. Create new style system alongside existing code
2. Gradually migrate components one by one
3. Maintain backward compatibility during transition
4. Remove old patterns once migration is complete
5. Update documentation and examples

---

**Next Steps:** Start with Phase 1, Task 1.1 - creating the style module and extracting existing styles.