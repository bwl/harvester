# ENGINE.md

Packages
- pkg/engine: core game types and systems
- pkg/data: enums and constants

GameState
- Seeded RNG: *rand.Rand stored in GameState
- Fields: Tick, Epoch, Map, Player, Log (ring), Events

Map/Expansion
- Map has Width, Height, Tiles [][]Tile. Expand(growth) appends rim rows/cols; generates new tiles via Generator using RNG and Epoch weights.

Systems
- ExpansionSystem: on each player action increases size
- MovementSystem: validates movement vs drive range and terrain
- HarvestSystem: collects resources from Galaxy tiles
- HazardSystem: updates effects (pull, slow)
- UpgradeSystem: applies upgrades affecting rules

APIs
- engine.New(seed int64) *GameState
- (*GameState) Step(action Action) Result // processes one action and advances systems

Determinism
- All randomness via GameState.RNG; tests use fixed seeds.
