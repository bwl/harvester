# Harvester 🚀

A terminal-based space exploration roguelike where you pilot a rocket ship to explore mysterious planets in a multi-layered universe.

## 🎮 Game Overview

**Harvester** is a multi-layered space exploration roguelike built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea). Players begin in space with a rocket ship and must strategically choose which planets to explore, knowing that once they land, they're committed to that world until they complete the necessary quests to re-enable space travel.

### Core Gameplay

- **Space Navigation**: Pilot your rocket ship through space with limited fuel
- **Strategic Planet Choice**: Choose from 3 unique planets, each with distinct biomes and challenges  
- **Deep Exploration**: Explore planets hundreds of levels deep with Angband-style progression
- **Quest-Gated Progression**: Complete planet-specific objectives to escape and continue your journey
- **Cross-Planet Synergies**: Knowledge and items from previous planets help on new worlds

## 🚀 Quick Start

### Prerequisites

- Go 1.23.0 or later
- Terminal with color support

### Installation & Running

```bash
# Clone the repository
git clone https://github.com/bwl/harvester.git
cd harvester

# Build and run
make build
make run

# Or run directly
go run ./cmd/game
```

### Development Commands

```bash
# Build
make build          # or: go build ./...

# Testing
make test           # or: go test ./...
make clean          # or: go clean -testcache

# Code Quality
make fmt            # or: go fmt ./...
make vet            # or: go vet ./...
make lint           # or: golangci-lint run (requires golangci-lint)
```

## 🏗️ Architecture

Harvester is built with a modern ECS (Entity-Component-System) architecture:

- **Entry Point**: `cmd/game/main.go` - Initializes the Bubble Tea TUI
- **Game Engine**: `pkg/engine/` - Core game systems and universe expansion
- **ECS Framework**: `pkg/ecs/` - Entity-Component-System implementation
- **UI Layer**: `internal/ui/` - Bubble Tea interface with lipgloss styling
- **Game Systems**: `pkg/systems/` - Movement, rendering, combat, and more
- **Components**: `pkg/components/` - Position, health, player, tiles, etc.

### Key Features

- **Deterministic World Generation**: Seeded RNG ensures reproducible universes
- **Layered Rendering**: Alpha blending system for complex visual effects
- **Performance Optimized**: Efficient ECS queries and component storage
- **Extensible Design**: Easy to add new components, systems, and mechanics

## 🎯 Game Mechanics

### Game Loop Structure

#### **1. Space Layer - Navigation & Choice**
- Start in a rocket ship with limited fuel
- **3 planets** available for exploration, each with unique characteristics
- **Fuel management** creates strategic tension
- Planets show basic information (biome type, difficulty hints)

#### **2. Planet Landing - Commitment Point**
- Landing is **irreversible** - you're committed to this planet
- Rocket ship becomes inoperable until quest conditions are met
- Planet surface reveals the chosen biome and exploration area

#### **3. Deep Exploration - The Real Game**
- Each planet can be **hundreds of levels deep**
- **Angband-style depth progression** with increasing difficulty and rewards
- **Biome-specific systems** create unique gameplay on each planet
- Quest objectives scattered throughout the depths

#### **4. Quest Gates - Escape Condition**
- Each planet has **unique escape requirements** 
- Examples: Craft heat shield, repair technology, tame wildlife
- Completing requirements **re-enables space travel**

### Planet Types & Biomes
- **Desert Worlds**: Heat management and water scarcity
- **Ice Worlds**: Cold survival and thermal regulation
- **Jungle Worlds**: Dense vegetation and exotic life forms
- **Volcanic Worlds**: Lava navigation and mineral extraction
- **Ancient Worlds**: Mysterious ruins and forgotten technology

## 🛠️ Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal User Interface framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling and layout
- [Harmonica](https://github.com/charmbracelet/harmonica) - Spring animations
- [Testify](https://github.com/stretchr/testify) - Testing utilities

## 📚 Documentation

Comprehensive documentation is available in the `docs/` directory:

- [**DESIGN.md**](docs/DESIGN.md) - Game design and mechanics
- [**ENGINE.md**](docs/ENGINE.md) - Engine architecture and APIs  
- [**UI.md**](docs/UI.md) - UI components and styling
- [**KEYMAP.md**](docs/KEYMAP.md) - Controls and key bindings
- [**ECS_PLAN.md**](docs/ECS_PLAN.md) - Entity-Component-System architecture
- [**RUNNING.md**](docs/RUNNING.md) - Build, run, and development guide

## 🎮 Controls

- **Arrow Keys / WASD**: Navigate your ship or character
- **Enter**: Interact/Select
- **Tab**: Menu navigation
- **Esc**: Back/Pause
- **Q**: Quit game
- **Save/Load**: Automatic saves with manual slot system

## 🤝 Contributing

This is a personal project, but feedback and suggestions are welcome! Check out the documentation in `docs/` to understand the codebase architecture.

## 📄 License

This project is open source. See the repository for license details.

---

*Explore strange new worlds, one planet at a time.* 🌌