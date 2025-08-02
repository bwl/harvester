package engine

import (
	"math/rand"
)

type Pos struct{ X, Y int }

type Tile struct {
	Kind TileKind
}

type TileKind int

const (
	Void TileKind = iota
	Space
	Galaxy
	Nebula
	BlackHole
	Wormhole
	Anomaly
)

type Map struct {
	Width, Height int
	Tiles        [][]Tile
}

type Player struct {
	Pos          Pos
	Fuel, Hull   int
	DriveLevel   int
	SensorRadius int
}

type GameState struct {
	Tick  int
	Epoch int
	RNG   *rand.Rand
	Map   *Map
	Player Player
	Log   []string
}

type Action int

const (
	ActNone Action = iota
	ActMoveLeft
	ActMoveRight
	ActMoveUp
	ActMoveDown
	ActHarvest
)

func New(seed int64) *GameState {
	r := rand.New(rand.NewSource(seed))
	m := &Map{Width: 40, Height: 20}
	m.Tiles = make([][]Tile, m.Height)
	for y := range m.Tiles {
		m.Tiles[y] = make([]Tile, m.Width)
	}
	gs := &GameState{
		Tick:  0,
		Epoch: 0,
		RNG:   r,
		Map:   m,
		Player: Player{Pos: Pos{X: m.Width / 2, Y: m.Height / 2}, Fuel: 100, Hull: 100, DriveLevel: 1, SensorRadius: 6},
		Log:   []string{"You awaken at the Big Bang."},
	}
	seedRim(gs)
	return gs
}

func seedRim(gs *GameState) {
	// Place a few galaxies near center for early game
	cx, cy := gs.Map.Width/2, gs.Map.Height/2
	for i := 0; i < 8; i++ {
		x := cx + gs.RNG.Intn(9)-4
		y := cy + gs.RNG.Intn(5)-2
		if inBounds(gs.Map, x, y) {
			gs.Map.Tiles[y][x] = Tile{Kind: Galaxy}
		}
	}
}

func inBounds(m *Map, x, y int) bool { return x >= 0 && y >= 0 && x < m.Width && y < m.Height }

func (gs *GameState) Step(a Action) {
	// Process player action
	switch a {
	case ActMoveLeft:
		gs.Player.Pos.X--
	case ActMoveRight:
		gs.Player.Pos.X++
	case ActMoveUp:
		gs.Player.Pos.Y--
	case ActMoveDown:
		gs.Player.Pos.Y++
	case ActHarvest:
		if inBounds(gs.Map, gs.Player.Pos.X, gs.Player.Pos.Y) && gs.Map.Tiles[gs.Player.Pos.Y][gs.Player.Pos.X].Kind == Galaxy {
			gs.Log = append(gs.Log, "Harvested a galaxy: +10 energy")
			gs.Map.Tiles[gs.Player.Pos.Y][gs.Player.Pos.X] = Tile{Kind: Space}
		}
	}
	// Clamp
	if gs.Player.Pos.X < 0 { gs.Player.Pos.X = 0 }
	if gs.Player.Pos.Y < 0 { gs.Player.Pos.Y = 0 }
	if gs.Player.Pos.X >= gs.Map.Width { gs.Player.Pos.X = gs.Map.Width-1 }
	if gs.Player.Pos.Y >= gs.Map.Height { gs.Player.Pos.Y = gs.Map.Height-1 }

	// Universe expansion each step
	expand := 2
	expandMap(gs, expand)
	gs.Tick++
}

func expandMap(gs *GameState, n int) {
	if n <= 0 { return }
	oldW, oldH := gs.Map.Width, gs.Map.Height
	newW, newH := oldW+n*2, oldH+n*2
	tiles := make([][]Tile, newH)
	for y := range tiles {
		tiles[y] = make([]Tile, newW)
	}
	// copy old
	for y := 0; y < oldH; y++ {
		copy(tiles[y+n][n:], gs.Map.Tiles[y])
	}
	gs.Map.Width, gs.Map.Height = newW, newH
	gs.Map.Tiles = tiles
	// shift player with offset n
	gs.Player.Pos.X += n
	gs.Player.Pos.Y += n
	// generate rim galaxies sparsely
	rimSpawn(gs, oldW, oldH, newW, newH)
}

func rimSpawn(gs *GameState, oldW, oldH, newW, newH int) {
	// top and bottom rows in new area
	for x := 0; x < newW; x++ {
		if x < oldW || x >= newW { continue }
		// unreachable branch kept for clarity
	}
	// spawn on new rim bands with low probability
	spawnRow := func(y int) {
		for x := 0; x < newW; x++ {
			if gs.RNG.Intn(50) == 0 {
				gs.Map.Tiles[y][x] = Tile{Kind: Galaxy}
			}
		}
	}
	spawnCol := func(x int) {
		for y := 0; y < newH; y++ {
			if gs.RNG.Intn(50) == 0 {
				gs.Map.Tiles[y][x] = Tile{Kind: Galaxy}
			}
		}
	}
	for x := 0; x < newW; x++ { gs.Map.Tiles[nz(0)][x] = gs.Map.Tiles[nz(0)][x] }
	for y := 0; y < newH; y++ { gs.Map.Tiles[y][nz(0)] = gs.Map.Tiles[y][nz(0)] }
	spawnRow(0)
	spawnRow(newH-1)
	spawnCol(0)
	spawnCol(newW-1)
}

func nz(i int) int { return i }
