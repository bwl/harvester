# NEXT-STEPS.md

Vision
- Big Bang roguelike: each move expands the universe; harvest galaxies; upgrade propulsion to keep up with expansion. Charm (Bubble Tea, Lip Gloss) for TUI.

Core pillars
- Responsive TUI at 60fps max without flicker
- Deterministic game logic with seeded RNG
- Simple, testable ECS-ish architecture
- Accessibility: configurable colors/keys

High-level architecture
- Engine: expansion system, RNG, rim generator, FOV, pathfinding
- UI: Bubble Tea model/update/view, components (status, map, log, upgrades)
- Styling: Lip Gloss theme system, high-contrast fallback
- Data: tiles/resources/upgrades as Go structs, JSON-compatible

Key systems
- Map/tiles: grid with materials; rooms+corridors; cellular automata smoothing; seeded RNG
- FOV/lighting: shadow casting (Bresenham/Permissive) with visibility cache
- Actors: player, NPCs; components: Position, Stats, AI, Inventory
- Input: vim/arrows, actions mapped via keymap; command queue
- Turns/AI: energy-based or simple alternate turns; pathfinding A* using walk-cost
- Combat: melee, damage types based on material interactions
- Log: ring buffer shown in UI panel

Bubble Tea integration
- Model holds GameState and UI state; messages for Tick, Key, Resize
- Use tea.Batch for coalescing updates; avoid heavy work in View
- Use alt screen and mouse optional; disable when not supported

Performance
- Render diff-friendly strings; precompute styled glyphs per material
- Avoid rebuilding large strings; reuse buffers

Testing
- Pure functions for gen/AI/pathfinding; golden tests for FOV; seed control via rand.Rand

Milestones
1) Core loop: input -> update -> render with expanding map per move
2) Rim generator v1: sparse galaxy spawns on expansion; harvest action
3) Movement + sensor-based FOV + basic log
4) Upgrades UI and drive/sensor upgrades affecting mechanics
5) Hazards (black holes, nebula) + warp ability
6) Resources economy and progression goals
7) Polish: themes, key remap, save/load (optional)

Tech debt guardrails
- No globals except configured RNG source
- Context-aware operations for future I/O
- Small packages: engine, ui, data, util

Open questions
- Turn model (energy vs strict turns)
- Map size and camera viewport
- Persistence scope

Immediate next steps
- Extract engine package (map, rng source, rect, types)
- Add internal/ui package with Bubble Tea model decomposition
- Add FOV and keymap scaffolding
- Introduce go test scaffolding and basic unit tests for RNG determinism
