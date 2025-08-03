# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build and Run
- Build: `go build ./...` or `make build`
- Run: `go run ./cmd/game` or `make run`
- Clean test cache: `go clean -testcache` or `make clean`

### Linting and Formatting
- Format code: `go fmt ./...` or `make fmt`
- Fix imports: `goimports -w .` or `make imports` (requires goimports)
- Vet: `go vet ./...` or `make vet`
- Static check: `staticcheck ./...` or `make staticcheck` (requires staticcheck)
- Lint: `golangci-lint run` or `make lint` (requires golangci-lint)

### Testing
- Run all tests: `go test ./...` or `make test`
- Run single package tests: `go test ./path/to/pkg -v`
- Run single test: `go test ./path/to/pkg -run ^TestName$ -v`
- Run benchmarks: `go test ./path/to/pkg -bench . -benchmem`

## Architecture

This is a Big Bang roguelike game built with Go, using Bubble Tea for the TUI interface.

### Core Structure
- **Entry point**: `cmd/game/main.go` - initializes game state and Bubble Tea program with alt screen
- **Game engine**: `pkg/engine/engine.go` - contains core game logic, universe expansion, player movement
- **UI layer**: `internal/ui/model.go` - Bubble Tea model handling input/rendering with lipgloss styling
- **Data types**: `pkg/data/types.go` - shared type definitions

### Key Concepts
- **Universe expansion**: Each player action expands the map by adding rim tiles with procedural generation
- **Game state**: Centralized in `engine.GameState` with tick-based progression
- **Rendering**: Uses lipgloss for styled terminal output with viewport centering on player
- **Seeded RNG**: Deterministic world generation using `rand.New(rand.NewSource(seed))`

### Module Structure
- `pkg/engine` - Core game systems (expansion, movement, harvest, map generation)
- `internal/ui` - Bubble Tea UI components and styling
- `pkg/data` - Type definitions for tiles, upgrades, resources
- `cmd/game` - Main application entry point

### ECS Migration Plan (docs/ECS_PLAN.md)
The codebase is planned to migrate to an Entity-Component-System (ECS) architecture:
- **Phase 1**: Core ECS with World, Entity registry, component stores, queries, system scheduler
- **Phase 2**: Migrate existing engine state to ECS components and systems
- **Phase 3**: Add tests, benchmarks, and quality improvements
- **Planned packages**: `pkg/ecs/`, `pkg/components/`, `pkg/systems/`
- **Integration**: Keep Bubble Tea as controller forwarding input to ECS and rendering projections

## Dependencies

Primary dependencies:
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Terminal styling

Optional development tools (check availability before use):
- `goimports` - for import formatting
- `staticcheck` - for static analysis  
- `golangci-lint` - for comprehensive linting

## Code Style

Follow existing patterns:
- Use explicit types for exported APIs
- Keep functions small with early returns
- Handle errors with `fmt.Errorf("...: %w", err)`
- Group imports: standard library, third-party, local
- Use concrete types locally, avoid `interface{}`
- Receivers use short names (gs, m, p)