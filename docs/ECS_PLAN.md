# ECS Implementation Plan

## Goals
- Introduce a lightweight, testable ECS to structure gameplay logic.
- Keep UI (Bubble Tea) decoupled from ECS runtime.
- Enable incremental migration from current engine to ECS.

## Scope (Phase 1)
- Core ECS with entity registry, generic component stores, queries, and system scheduler.
- Essential components (Position, Velocity, Renderable, Health, Input).
- Minimal systems (Movement, Rendering adapter, Input adapter).

## Package Layout
- pkg/ecs/
  - world.go (World, Schedule)
  - entity.go (Entity IDs)
  - store.go (component storage, sparse-set/map)
  - query.go (iterators over component sets)
  - event.go (optional simple bus)
  - system.go (System interface)
- pkg/components/
  - position.go, velocity.go, renderable.go, health.go, input.go
- pkg/systems/
  - movement.go, render.go, input.go
- cmd/game/
  - wire ECS world into existing main and UI loop

## Data Model
- type Entity = uint64 (or ecs.Entity struct wrapping id)
- Component: plain Go structs; one type per component; no interfaces.
- Storage: map[Entity]T for simplicity initially; hide behind generics store[T any]. Later upgrade to sparse-set.

## Public API
- World:
  - Create() Entity
  - Destroy(e Entity)
  - Add[T any](e Entity, c T)
  - Remove[T any](e Entity)
  - Get[T any](e Entity) (T, bool)
  - Query[T1, T2, ...]() iterator
- System:
  - Update(dt float64, w *World)
- Scheduler:
  - Fixed order list with priorities: Input -> AI -> Physics -> Collision -> Combat -> Rendering

## Cross-cutting Concerns

### Error Handling
- Get returns (zero, false); Add/Remove return error; Destroy is idempotent.
- Debug mode: optional asserts for invalid entity/component operations.

### Deterministic Execution
- Fixed scheduler order; single-threaded Update; fixed dt accumulator; RNG owned by World and explicitly seeded.

### Memory Management
- Entity uses (index, generation); free-list reuse increments generation; Destroy purges all component stores; stores avoid retaining tombstones.

### Migration Compatibility
- Versioned save format; engine↔ECS adapters; during transition, persist both representations and prefer ECS when present.

### Performance Targets (Phase 1)
- 10k entities, 3 comps each; Query2 full-iterate ≤1ms/frame; Movement ≤0.5ms/frame; 60 FPS end-to-end on laptop.

### Testing Strategy
- Property tests for Create/Destroy/Add/Remove invariants and determinism with fixed RNG+dt.
- Golden snapshot tests for system state transitions; fuzz Query2/3; benches for queries and Movement.

## Phase-by-Phase Plan

### Phase 1: Core ECS
1. pkg/ecs/world.go: implement World with:
   - entity counter, free-list
   - component stores registry keyed by reflect.Type
   - Add/Get/Remove using a generic store[T]
2. pkg/ecs/store.go: implement store[T] with map[Entity]T and an index set.
3. pkg/ecs/query.go: implement Query2[T1,T2], Query3[…]; yield tuples and entity.
4. pkg/ecs/system.go: define System interface and basic Scheduler with ordered list.
5. pkg/components: add Position, Velocity, Renderable, Health, Input components.
6. pkg/systems/movement.go: integrate Position+Velocity update by dt.
7. pkg/systems/render.go: expose render data; integrate with UI via adapter function that reads Renderable + Position.
8. pkg/systems/input.go: translate Bubble Tea key msgs into Input components.
9. Wire into cmd/game: create world, register systems, drive scheduler each tick from update loop.

### Phase 2: Migration of Existing Engine
1. Identify current engine state in pkg/engine and map to components.
2. Move movement/update logic into systems; keep engine as façade delegating to ECS.
3. Replace direct state mutations with World.Add/Remove/Get.
4. Update UI model to read-only projection from ECS (no writes in view layer).

### Phase 3: Quality & Tooling
1. Tests: unit tests for store, queries, and movement system.
2. Benchmarks: Query2/Query3 iteration microbenchmarks.
3. Lint/typecheck in CI; ensure deterministic entity IDs in tests.

### Phase 4: Enhancements (Optional)
1. Sparse-set storage for O(1) packed iteration.
2. Event bus for decoupled interactions (Damage, Spawn, Despawn).
3. Prefabs/factories for entity construction.
4. Save/load: serialize components.
5. System dependencies and conditional execution.

## Integration with Bubble Tea
- Keep tea.Model as controller to forward input to ECS and render projections from ECS state.
- Rendering system produces a slice of drawables; UI renders via lipgloss.

## Risks & Mitigations
- Complexity creep: limit Phase 1 to simple map-backed stores.
- Performance: add benchmarks before optimizing to sparse-set.
- Coupling: maintain clear boundaries between ECS and UI/engine.

## Definition of Done (Phase 1)
- make run builds and runs using ECS-backed loop.
- Movement and basic rendering work via ECS systems with deterministic results for fixed seed+dt.
- go test ./... passes with core ECS tests, property tests, and benches compile; basic perf targets met.
