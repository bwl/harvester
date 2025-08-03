# Refactor Plan

## Objectives
- Fix movement/input bugs and remove hardcoded entity IDs
- Unify component definitions and system boundaries
- Stabilize layered scheduling; 20 FPS engine, 60 FPS UI
- Improve file structure and introduce missing core abstractions

## Guiding Principles
- Components are pure data in pkg/components
- Systems contain logic in pkg/systems and operate via queries, not IDs
- Layers are enforced by the scheduler, not per-system guards
- UI remains a thin adapter over ECS (render + input translation)

## Work Breakdown

### 1) Player Identity and Queries (P0)
- Add components/player.go: `type Player struct{}`
- Tag player entity in UI model creation
- Replace hardcoded entity IDs with ECS queries across systems
  - Introduce query helpers where missing (Tuple2/3 over Input, Position, Velocity, Player)
- Add minimal entity service: `pkg/ecs/entityref.go` with `GetSingleton[T any]` pattern if needed

### 2) Velocity Unification (P0)
- Remove systems.Velocity; use components.Velocity exclusively
- Migrate all movement systems to read/write components.Velocity
- Normalize field names to VX/VY (or X/Y) and update usages consistently

### 3) Component/System Boundary Cleanup (P0)
- Move FuelTank and any other data structs from systems to components
- Ensure no business state types remain in pkg/systems
- Add import checks via go vet/staticcheck to prevent cycles

### 4) Layered Scheduler Ownership (P1)
- Keep layer activation in ecs.SchedulerWithContext
- Remove manual layer checks from systems; assign systems to Space/Surface/Deep lists
- Add transition hooks: optional `OnEnter/OnExit(layer)` in a lightweight manager if required

### 5) Input Flow Consistency (P1)
- Normalize input: InputSystem sets components.Input; movement systems consume it
- Space: velocity-based integration; Surface: discrete tile steps using same Input component
- Provide small adapter in SurfaceMovement to quantize movement while still reading velocity/input

### 6) UI Slimming (P1)
- Extract world/bootstrap to `pkg/engine/bootstrap.go`
- Keep UI model focused on Bubble Tea messages, rendering, save/load triggers
- Add renderer interface in systems.Render/MapRender kept as systems but UI only reads outputs

### 7) File Organization (P1)
- pkg/components: position.go, velocity.go, input.go, health.go, renderable.go, gameplay.go, tile.go, stats.go, camera.go, player.go, worldinfo.go
- pkg/systems: input.go, movement.go (space/surface subfiles), camera.go, maprender.go, render.go, selection.go, combat.go, levels.go, tick.go, fuel.go
- Remove mixed types from space.go; split into focused files

### 8) Persistence and Saves (P2)
- Ensure new components are included in save/load snapshot
- Add migrations for moved/renamed components where necessary

### 9) Performance & Rates (P2)
- Confirm engine tick at 20 FPS with fixed dt (0.05)
- UI tick at 60 FPS; render reads latest system outputs without extra Updates
- Add simple profiler hooks around scheduler.Update when running with DEBUG env

## Milestones
- M1 (P0): Player tagging + query-based movement; velocity/type unification; movement bug fixed
- M2 (P1): Boundary cleanup; layered scheduler ownership; input flow stabilized; UI slimmed bootstrap
- M3 (P1): File reorg; systems split; build/tests green; saves compatible
- M4 (P2): Performance verification; basic profiling; docs updated

## Risk Mitigation
- Incremental PRs per milestone; compile after each step
- Run go fmt, vet, and go build per change; keep snapshots loadable
- Add minimal tests in pkg/testharness covering movement, layer transitions, and save/load
- Test movement in both space and surface layers after each P0 change
- Verify save/load works with Player component before proceeding to P1

## Acceptance Criteria
- No hardcoded entity IDs; player found via component/tag
- Single Velocity type used everywhere
- Systems free of component type definitions
- Scheduler controls layers; no per-system layer guards
- Game runs: engine 20 FPS, UI 60 FPS; movement works in space and surface
- Build, vet pass clean; saves load correctly
