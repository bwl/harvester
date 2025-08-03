# Layered Rendering Architecture Plan (MVP)

## Problem Analysis
Manual ANSI string splicing causes corruption. We need a glyph-based compositor with clean layering.

## Proposed Architecture

MVP targets Bubble Tea RPG needs with full WindowSize responsiveness, layers GAME + MENU only, Lip Gloss rasterization for borders/panels, alignment+offset positioning, no animations, full redraw per frame, and basic wide-rune support.

### Master View Controller
A pure `ViewRenderer` that composes layers into a terminal string. No external state mutations in Render(); Bubble Tea calls View() to get its output.

### Layer System (MVP)
1. GAME — terrain/map/entities
2. MENU — dialogs/menus/forms
(Defer UI/HUD and OVERLAY.)

### Interfaces
- RenderableContent: GetLayer, GetZ, GetPosition, GetBounds, GetGlyphs
- Position: Horizontal/Vertical alignment + OffsetX/OffsetY (percentages deferred)
- Glyph: rune + fg/bg + style (bold/italic/etc.)
- Bounds, Color, Style types

## Bubble Tea Integration
- Model.Update(): register/unregister content based on state
- Model.View(): vr.Render()
- tea.WindowSizeMsg: vr.SetDimensions(w,h); minimum recommended 80x24; clip/crop or center on smaller

## Lip Gloss Integration (MVP)
- Use lipgloss for borders/panels/basic styling only
- Add rasterizer: LipGlossString → [][]Glyph to avoid ANSI corruption

## Positioning (MVP)
- Alignment + offsets only (Left/Center/Right; Top/Center/Bottom)
- Add placeholders for percentage-based positioning in types/APIs (fields unused in MVP) to avoid refactor later.

## Implementation Plan

### Phase 1: Core
- pkg/rendering/interfaces.go
  - Layer (Game, Menu; defer UI, Overlay)
  - Position (alignment+offset only)
  - Bounds, Glyph, Color, Style
  - RenderableContent
  - [ ] TODO: Add PercentX/PercentY fields (unused in MVP) for future percentage positioning
- pkg/rendering/glyph_matrix.go
  - GlyphMatrix, NewGlyphMatrix, Clear, Set/Get, InBounds
  - ToTerminalString (handles wide runes; combining deferred; flag combining marks as high-priority follow-up for CJK/emoji)
  - [ ] TODO: Combine mark handling (post-MVP)
- pkg/rendering/position.go
  - CalculatePosition for alignment+offset; bounds clipping
- pkg/rendering/view_renderer.go
  - width/height, layers map, glyph matrix
  - Register/Unregister, SetDimensions, Render (pure), composite
- pkg/rendering/lipgloss_raster.go
  - Convert lipgloss-rendered blocks to [][]Glyph
  - Define mapping: Lip Gloss styles/colors → Glyph.Style/Color (bold, italic, underline, dim, reverse)
  - [ ] TODO: Document exact style/color mapping table

### Lip Gloss → Glyph Mapping
- Styles:
  - lipgloss.Bold() → StyleBold
  - lipgloss.Italic() → StyleItalic
  - lipgloss.Underline() → StyleUnderline
  - lipgloss.Faint()/Dim() → StyleDim
  - lipgloss.Reverse() → StyleReverse
- Colors:
  - Foreground: map to Glyph.Foreground as 24-bit RGB when available; zero Color = no color
  - Background: map to Glyph.Background as 24-bit RGB when available; zero Color = no color
- Rendering:
  - matrixToString emits ANSI only on change; resets at line start and end; zero Color resets to clear styles
  - RenderLipglossString strips ANSI from input strings before rasterization

### Phase 2: StartScreen Migration
- internal/ui/screens/terrain_content.go (LayerGame): wrap map render → glyphs
- internal/ui/screens/menu_content.go (LayerMenu): menu via lipgloss → raster → glyphs
- startscreen: hold ViewRenderer; Update() manages registrations; View() uses vr.Render(); handle WindowSizeMsg
- Tests: terrain shows; menu centers; resize works; golden snapshots

### Phase 3: Post-MVP Enhancements
- Z-index: global ordering via GetZ(); Layer is semantic only.
- Conventions:
  - ZBackground=0, ZContent=100, ZUI=120, ZHUD=130, ZMenu=150, ZPattern=800, ZFrame=1000
  - UI elements start at >=100; background strictly 0; frame highest
  - Equal Z preserves registration order (stable sort)
- Add UI/HUD layer via lipgloss rasterizer
- Simple animations (blink/fade) once static composition is stable
- Clipping/viewport and Z ordering within layers as needed

### Phase 4: Optimization (MVP+1)
- Add basic frame timing instrumentation (optional log flag) to measure redraw costs
- [ ] TODO: Frame timing instrumentation hook
- Dirty region tracking; partial redraws
- Glyph matrix pooling; faster glyph->string conversion
- Benchmarks

### Testing
- Unit tests: position calc, glyph matrix, renderer
- Golden snapshot tests for View() output (include wide-rune cases)
- Parameterized resize tests across widths/heights to validate alignment and clipping
- [ ] TODO: Table-driven resize tests for alignment/centering/clipping

## File Structure
pkg/rendering/
- interfaces.go
- colors.go
- glyph_matrix.go
- position.go
- view_renderer.go
- lipgloss_raster.go

internal/ui/screens/
- terrain_content.go
- menu_content.go
# ui_content.go, animation_content.go deferred
