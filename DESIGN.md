# DESIGN.md

Title: Big Bang Roguelike (Charm TUI)

Premise
- You appear at the Big Bang. Each move expands the universe (map). Harvest galaxies to progress. Expansion increases distances; upgrade propulsion to keep pace.

Core loop
- Input (move/action) -> Tick: expand universe, advance entities, spawn galaxies -> Resolve interactions (harvest, combat hazards) -> Render.

Objectives/Failure
- Win: reach target harvested energy/milestones across epochs.
- Lose: run out of fuel/hull or get trapped by expansion hazards.

World model
- Infinite expanding grid. Visible viewport centers on player.
- Tiles: Void, Space, Galaxy, Nebula, BlackHole, Wormhole, Anomaly.
- Epochs: parameters change over time (spawn rates, hazards).

Player systems
- Resources: fuel, hull, energy, data.
- Upgrades: Drive (range/AP), Sensors (FOV), Cargo (harvest yield), Shield (damage), Warp (blink), Autopilot (pathing efficiency).
- Actions: move, harvest, warp, scan, craft/upgrade.

Mechanics
- Expansion: each player action increases MapWidth/Height by growth_rate; new rim tiles populated by Procedural Gen with seeded RNG.
- Movement AP cost scales with distance; Drive level reduces cost or increases max step.
- Harvest: stand on Galaxy; gain resources proportional to size and tools; depletes node.
- Hazards: BlackHole pulls; Nebula slows; Anomaly events; Wormholes teleport.

Procedural generation
- Seeded rand.Rand; epoch-driven weights. Galaxy clustering via Poisson disk or jittered grid at rim bands.

FOV/Visibility
- Sensor-based radius; upgrades increase radius; fog-of-war kept.

UI (Bubble Tea + Lip Gloss)
- Layout: top status (epoch, resources, upgrades), center map viewport, right log, bottom actions/help.
- Components: MapView, StatusBar, LogView, Tooltip/Help, Modal (upgrade menu).
- Keymap: hjkl/arrows move; g harvest; w warp; s scan; u upgrades; q quit.

Architecture
- pkg/engine: GameState, Map, Gen, RNG, Systems (Expansion, Movement, Harvest, Hazards, Upgrades).
- internal/ui: model, components, styles, keymap.
- cmd/game: main entry.

Data
- Define types in pkg/data: TileKind, Resource, UpgradeKind; JSON tags for potential save.

Performance
- Only render viewport; precompute styled glyphs per TileKind; minimize allocations.

Testing
- Deterministic tests for gen, expansion, harvest yields, hazard effects; golden tests for FOV.

Extensibility
- Systems pattern; event bus in engine for UI messages.
