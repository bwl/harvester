package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	MapWidth  = 100
	MapHeight = 40

	RoomMinSize = 4
	RoomMaxSize = 15

	MaxRooms = 20

	// Materials
	MaterialRock  = "Rock"
	MaterialSoil  = "Soil"
	MaterialSand  = "Sand"
	MaterialWater = "Water"
	MaterialLava  = "Lava"
	MaterialGrass = "Grass"
	MaterialIce   = "Ice"
	MaterialMetal = "Metal"
)

type Tile struct {
	Type      TileType
	Material  string
	Energy    float64
	IsVisible bool
}

type TileType int

const (
	Wall TileType = iota
	Floor
)

type Rect struct {
	X, Y, W, H int
}

type Map struct {
	Tiles [][]Tile
}

func NewMap() *Map {
	m := &Map{
		Tiles: make([][]Tile, MapHeight),
	}
	for i := range m.Tiles {
		m.Tiles[i] = make([]Tile, MapWidth)
		for j := range m.Tiles[i] {
			m.Tiles[i][j] = Tile{
				Type:      Wall,
				Material:  MaterialRock,
				Energy:    rand.Float64(),
				IsVisible: false,
			}
		}
	}
	return m
}

func (m *Map) CreateRoom(room Rect) {
	materials := []string{
		MaterialRock, MaterialSoil, MaterialSand,
		MaterialWater, MaterialLava, MaterialGrass,
		MaterialIce, MaterialMetal,
	}

	// Choose a random material for the entire room
	material := materials[rand.Intn(len(materials))]

	for y := room.Y; y < room.Y+room.H; y++ {
		for x := room.X; x < room.X+room.W; x++ {
			m.Tiles[y][x] = Tile{
				Type:      Floor,
				Material:  material,
				Energy:    rand.NormFloat64(),
				IsVisible: true,
			}
		}
	}
}

func (m *Map) CreateHTunnel(x1, x2, y int) {
	for x := x1; x <= x2; x++ {
		m.Tiles[y][x] = Tile{
			Type:      Floor,
			Material:  MaterialSand,
			Energy:    rand.NormFloat64(),
			IsVisible: true,
		}
	}
}

func (m *Map) CreateVTunnel(y1, y2, x int) {
	for y := y1; y <= y2; y++ {
		m.Tiles[y][x] = Tile{
			Type:      Floor,
			Material:  MaterialSand,
			Energy:    rand.NormFloat64(),
			IsVisible: true,
		}
	}
}

func (m *Map) Generate() []Rect {
	rooms := []Rect{}

	for i := 0; i < MaxRooms; i++ {
		w := rand.Intn(RoomMaxSize-RoomMinSize) + RoomMinSize
		h := rand.Intn(RoomMaxSize-RoomMinSize) + RoomMinSize
		x := rand.Intn(MapWidth - w - 1)
		y := rand.Intn(MapHeight - h - 1)

		newRoom := Rect{X: x, Y: y, W: w, H: h}

		overlapping := false
		for _, otherRoom := range rooms {
			if newRoom.X <= otherRoom.X+otherRoom.W && newRoom.X+newRoom.W >= otherRoom.X &&
				newRoom.Y <= otherRoom.Y+otherRoom.H && newRoom.Y+newRoom.H >= otherRoom.Y {
				overlapping = true
				break
			}
		}

		if !overlapping {
			m.CreateRoom(newRoom)
			if len(rooms) != 0 {
				prevRoom := rooms[len(rooms)-1]
				if rand.Intn(2) == 0 {
					m.CreateHTunnel(prevRoom.X+prevRoom.W/2, newRoom.X+newRoom.W/2, prevRoom.Y+prevRoom.H/2)
					m.CreateVTunnel(prevRoom.Y+prevRoom.H/2, newRoom.Y+newRoom.H/2, newRoom.X+newRoom.W/2)
				} else {
					m.CreateVTunnel(prevRoom.Y+prevRoom.H/2, newRoom.Y+newRoom.H/2, prevRoom.X+prevRoom.W/2)
					m.CreateHTunnel(prevRoom.X+prevRoom.W/2, newRoom.X+newRoom.W/2, newRoom.Y+newRoom.H/2)
				}
			}
			rooms = append(rooms, newRoom)
		}
	}

	m.applyGaussianNoise()
	return rooms
}

func (m *Map) applyGaussianNoise() {
	for y := 0; y < MapHeight; y++ {
		for x := 0; x < MapWidth; x++ {
			if m.Tiles[y][x].Type == Wall {
				m.Tiles[y][x].Energy += rand.NormFloat64()
				if m.Tiles[y][x].Energy < -1.5 {
					m.Tiles[y][x].Material = MaterialRock
				} else if m.Tiles[y][x].Energy < -1 {
					m.Tiles[y][x].Material = MaterialSoil
				} else if m.Tiles[y][x].Energy < -0.5 {
					m.Tiles[y][x].Material = MaterialSand
				} else if m.Tiles[y][x].Energy < 0 {
					m.Tiles[y][x].Material = MaterialWater
				} else if m.Tiles[y][x].Energy < 0.5 {
					m.Tiles[y][x].Material = MaterialLava
				} else if m.Tiles[y][x].Energy < 1 {
					m.Tiles[y][x].Material = MaterialGrass
				} else if m.Tiles[y][x].Energy < 1.5 {
					m.Tiles[y][x].Material = MaterialIce
				} else {
					m.Tiles[y][x].Material = MaterialMetal
				}
			}
		}
	}
}

type Character struct {
	Name   string
	Health int
	X      int
	Y      int
}

type model struct {
	character Character
	gameMap   *Map
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "h", "left":
			if m.character.X > 0 && m.gameMap.Tiles[m.character.Y][m.character.X-1].Type == Floor {
				m.character.X--
			}
		case "j", "down":
			if m.character.Y < MapHeight-1 && m.gameMap.Tiles[m.character.Y+1][m.character.X].Type == Floor {
				m.character.Y++
			}
		case "k", "up":
			if m.character.Y > 0 && m.gameMap.Tiles[m.character.Y-1][m.character.X].Type == Floor {
				m.character.Y--
			}
		case "l", "right":
			if m.character.X < MapWidth-1 && m.gameMap.Tiles[m.character.Y][m.character.X+1].Type == Floor {
				m.character.X++
			}
		}
	}
	return m, nil
}

var (
	rockStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render
	soilStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("94")).Render
	sandStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("222")).Render
	waterStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("21")).Render
	lavaStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Render
	grassStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Render
	iceStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render
	metalStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Render
)

func (m model) View() string {
	grid := ""
	for y := 0; y < MapHeight; y++ {
		for x := 0; x < MapWidth; x++ {
			if m.character.X == x && m.character.Y == y {
				grid += "@"
			} else {
				switch m.gameMap.Tiles[y][x].Type {
				case Floor:
					switch m.gameMap.Tiles[y][x].Material {
					case MaterialRock:
						grid += rockStyle(".")
					case MaterialSoil:
						grid += soilStyle(".")
					case MaterialSand:
						grid += sandStyle(".")
					case MaterialWater:
						grid += waterStyle(".")
					case MaterialLava:
						grid += lavaStyle(".")
					case MaterialGrass:
						grid += grassStyle(".")
					case MaterialIce:
						grid += iceStyle(".")
					case MaterialMetal:
						grid += metalStyle(".")
					default:
						grid += "."
					}
				case Wall:
					grid += " " // Make walls appear as blank spaces
				}
			}
		}
		grid += "\n"
	}
	// Get the tile the character is currently on
	currentTile := m.gameMap.Tiles[m.character.Y][m.character.X]

	return fmt.Sprintf(
		"Character: %s\nHealth: %d\nPosition: (%d, %d)\nMaterial: %s\nEnergy: %.2f\n\n%s\nUse h/j/k/l or arrow keys to move. Press 'q' to quit.\n",
		m.character.Name, m.character.Health, m.character.X, m.character.Y, currentTile.Material, currentTile.Energy, grid,
	)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	character := Character{
		Name:   "Hero",
		Health: 100,
	}
	gameMap := NewMap()
	rooms := gameMap.Generate()

	if len(rooms) > 0 {
		startRoom := rooms[rand.Intn(len(rooms))]
		character.X = startRoom.X + startRoom.W/2
		character.Y = startRoom.Y + startRoom.H/2
	}

	p := tea.NewProgram(model{character: character, gameMap: gameMap})
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
