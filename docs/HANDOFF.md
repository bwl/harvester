# Project Handoff

Overview
- Language/stack: Go 1.23+, module "bubbleRouge". TUI game using Bubble Tea + Lip Gloss. ECS-based architecture.
- Entry points: cmd/game (interactive), cmd/sim (headless sim harness).
- Core ECS: pkg/ecs (World, stores, queries, scheduler, save/load), pkg/components (component types), pkg/systems (systems for input, movement, camera, render, tick, harvest, combat).
- UI: internal/ui/model.go wires systems, manages input, save/load.
- Tests: pkg/testharness contains controller harness and tests (movement, harvest, save/load, fuzz, concurrency).

ECS Save/Load Status
- Snapshot format (pkg/ecs/serialize.go):
  - Version (int) with migration hooks
  - Seed (int64) for RNG determinism
  - Allocator state: next, free entity list
  - Components map keyed by Go type name -> map[Entity]json.RawMessage
- Persisted components: Position, Velocity, Camera, PlayerStats, WorldInfo, Input, Inventory, Resource, Tile, Renderable, Health.
- Inventory post-unmarshal Ensure() added to guarantee map initialization.
- Encode/Decode helpers (pkg/ecs/compress.go): optional gzip compression and password-based CTR encryption (AES-CTR with SHA-256 key).
- Determinism: stores are cleared before load; allocator and seed restored; Save→Load→Save equivalence fuzz test in place.

UI Save/Load
- Keys: ctrl+s save JSON; ctrl+o load JSON; ctrl+shift+s save compressed (.gz); ctrl+shift+o load compressed.
- Save slots: 1/2/3 to save to .saves/slotN.gz; Shift+1/2/3 to load from slot.

Testing
- Unit/integration:
  - Movement and camera snapshot sanity.
  - Harvest resource -> inventory.
  - SaveLoad roundtrips for baseline components.
  - End-to-end controller Save/Load parity.
  - Encryption/Compression encode/decode roundtrip.
  - Fuzz equivalence test for Save→Load→Save (Go fuzz).
  - Concurrency test for parallel save/load.

Design Notes
- World holds typed stores (generic store[T]) and uses reflect.Type keys.
- Query helpers (View2/Each) for basic joins.
- Scheduler orders systems deterministically.
- World RNG default seed 1; persisted via Snapshot.Seed.

Open Items / Next Steps
1) Documentation
   - Update docs/UI.md and RUNNING.md with save/load keys, slot behavior, and file formats.
   - Add docs for Snapshot versioning and migration policy with examples.
2) CLI/Config
   - Flags/env vars for save directory, slot count, and password; surface in cmd/game and cmd/sim.
   - Add command to list slots with timestamps and metadata.
3) Persistence Extensions
   - As new components are added, extend Save/Load registrations and tests.
   - Consider binary encoder for performance (e.g., msgpack/cbor) behind an option.
   - Add snapshot checksums and basic integrity validation.
4) Testing & CI
   - Wire Go fuzzing in CI (nightly or bounded runs).
   - Add more property tests (random entity sets, removal/recreation, sparse stores).
   - UI integration tests for save slot keys (via Bubble Tea test harness if added).
5) Reliability & Migration
   - Implement snapshot migrations for Version < current (example: seed introduction).
   - Robust error handling on load; user feedback in UI (toasts/log).
   - Corruption handling: attempt partial recover or safe failure.
6) Performance
   - Benchmark Save/Load sizes and timings with/without compression.
   - Optimize store iteration and JSON allocation; consider pooling.
7) UX Improvements
   - On-screen hints for save/load keys; confirm messages.
   - Show slot listing with last-saved time and location in UI side panel.
8) Code Health
   - Add staticcheck/golangci-lint to dev workflow; go mod tidy.
   - Record standard commands in CRUSH.md updates if missing.

Quick Start Commands
- Build: go build ./...
- Run game: go run ./cmd/game
- Run tests: go test ./...
- Format/vet: go fmt ./... && go vet ./...
- Fuzz locally (Go >=1.18): go test ./pkg/testharness -run=^$ -fuzz=Fuzz -fuzztime=10s

Notes
- Secrets: encryption uses password-derived key; do not log secrets or store plaintext passwords.
- Determinism: keep scheduler order and RNG seed handling stable to maintain reproducible saves.
