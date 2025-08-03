# Test Harness Plan

## Goals
- Deterministic, headless simulation of game logic for CI and agents.
- Scripted inputs -> ECS world ticks -> JSON snapshots to assert.
- No Bubble Tea dependency during tests.

## Architecture
- pkg/testharness:
  - Controller struct wraps ecs.World and ecs.Scheduler, constructs minimal game world (player, camera, map, systems).
  - Methods:
    - InjectKey(key string): converts to Input component for the player.
    - Tick(n int, dt float64): runs scheduler n times with fixed dt.
    - Snapshot() ([]byte, error): returns canonical JSON of selected ECS components.
  - Options:
    - Seed int64, MapSize (w,h), PerfMode (skip render systems), IncludeComponents []string.
- cmd/sim:
  - Reads a JSON script: { seed, width, height, steps: [{key?:"left|right|up|down|g|...", ticks?:N}] }
  - Runs simulation and prints final snapshot to stdout.

## Determinism
- Single-threaded scheduler with fixed order.
- Fixed dt and RNG seeded per run.
- Avoid time.Now in harness; all randomness uses world RNG.

## Snapshot Format
- JSON with stable ordering:
  - player: { pos:{x,y}, stats:{fuel,hull,drive} }
  - entities: optional subsets (tiles/resources) summarized (counts per glyph in viewport) for compactness.
  - world: { tick, width, height }
- Avoid raw entity IDs where not needed; focus on observable state.

## Implementation Steps
1) Create pkg/testharness/controller.go
   - Build world: add systems InputSystem, Movement, CameraSystem(target=player), Tick, Harvest, Combat; optionally MapRender/Render behind flag.
   - Create player with Position, Velocity, Input, Renderable '@', PlayerStats; WorldInfo with map size.
   - Populate simple starfield/resources with seeded RNG.
   - Provide InjectKey -> systems.SetPlayerInput; Tick; Snapshot reading components.
2) Create cmd/sim/main.go
   - Parse script JSON from file or stdin; build Controller; loop steps; output Snapshot.
3) Tests
   - Test: moving right 3 ticks changes player X as expected.
   - Harvest test: place resource at player pos and press 'g' -> inventory increases.
   - Determinism test: same seed+script yields identical snapshot (golden file).
4) CI: add make targets `make sim` and `make test`; run sim with sample scripts in CI.

## Error Injection
- API to inject invalid states: missing components on controlled entity, out-of-bounds positions, negative stats; assert systems handle gracefully.
- Toggle to corrupt snapshots for load error paths.

## Save/Load Testing
- Harness exposes Save() and Load() mirroring UI bindings; tests: Save->Load->Snapshot equality (idempotent), and partial component restoration.

## Multi-Entity Scenarios
- Script ops: spawn {kind:"npc|enemy|resource", pos:{x,y}, comps:{...}}; enable adjacent combat tests and group behaviors.

## Performance Benchmarks
- Bench suite: Movement, Query2/3 iteration, Render adapter; target budgets (≤1ms/frame at 10k entities) and report.

## Golden File Management
- Store snapshots under testdata/*.golden.json; helper to update with UPDATE_GOLDEN=1; stable ordering and formatting.

## Script Validation
- JSON schema (docs/test_script.schema.json) enforced by harness; validate keys, types, ranges; clear errors.

## Extensions
- Property-based tests: random scripts produce consistent invariants (no negative HP, bounds not exceeded).
- Replay/record: capture inputs from a live session and replay in harness.
- Coverage of new systems as they ship.
