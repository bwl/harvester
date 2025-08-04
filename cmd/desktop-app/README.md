# Harvester Desktop Renderer

A native macOS desktop application for rendering the Harvester roguelike game using Swift and the game's Go engine.

## Architecture

```
Go Game Engine (bridge.go)
    ↓ C Bridge
Swift UI (SwiftUI + Canvas)
    ↓
Native macOS Window
```

## Features

- **Real-time rendering** of the game's glyph matrix
- **Native macOS performance** with SwiftUI Canvas
- **Full input handling** (WASD/Arrow keys for movement)
- **Rich color support** from the game's rendering system
- **Camera system** that follows the player
- **ECS integration** showing expanse tiles and entities

## Building and Running

### Three Versions Available

**1. macOS App Bundle (Recommended)** - Works with Command Line Tools only:
```bash
make build-app     # Create Harvester.app
open Harvester.app # Double-click to run
```

**2. Terminal Version** - Direct terminal execution:
```bash
make run           # Build and run in terminal
```

**3. Swift Version** - Requires full Xcode installation:
```bash  
make run-swift     # Build and run native macOS app
```

### Requirements

**App Bundle & Terminal Versions:**
- macOS 10.15+
- Command Line Tools for Xcode
- Go 1.19+
- GCC (included in Command Line Tools)

**Swift Version:**
- macOS 12.0+
- Full Xcode 14.0+ installation  
- Go 1.19+
- Swift 5.7+

### Installation

**Install to Applications folder:**
```bash
make install-app  # Builds and installs to /Applications/
```

**Manual installation:**
```bash
make build-app
cp -r Harvester.app /Applications/
```

### Clean
```bash
make clean
```

## Controls

- **W/↑**: Thrust forward
- **S/↓**: Brake
- **A/←**: Turn left  
- **D/→**: Turn right

## Components

### Go Bridge (`bridge.go`)
- **C structs**: `CGlyph`, `CGlyphMatrix` for Swift interop
- **Game state**: Manages player physics, ECS world, rendering
- **Exported functions**: `initGame`, `updateGame`, `getGlyphMatrix`

### Swift App
- **GameBridge**: Objective-C bridge to Go functions
- **GameView**: SwiftUI view with input handling and game loop
- **GlyphMatrixView**: Canvas-based renderer for glyphs

### Build System
- **Makefile**: Coordinates Go library and Swift app builds
- **Package.swift**: Swift package configuration with Go library linking

## Development

The app runs at 60 FPS, calling the Go engine's update and render functions each frame. The Go engine maintains full game state including:

- Player physics (position, velocity, rotation)
- ECS world with all entities and systems
- Camera system for viewport management
- Expanse tile generation and rendering

The Swift frontend focuses purely on:
- Input capture and forwarding to Go
- Efficient rendering of the glyph matrix
- Native macOS window management